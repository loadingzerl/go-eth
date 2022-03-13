package ut_tool

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/minio/minio-go/v6"
	"io"
	"io/ioutil"
)

const Endpoint string = "172.26.230.51:9000"

func PutJson(data []byte, height string) error {
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	minioClient, err := minio.New(Endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		fmt.Println(err)
		return err
	}
	bucketName := "lotus-state"
	location := "us-east-1"

	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(bucketName)
		if errBucketExists == nil && exists {
			//log.Printf("We already own %s\n", bucketName)
		} else {
			return err
		}
	}
	// Upload the zip file
	contentType := "application/octet-stream"
	param := Gzip(data)
	buf := bytes.NewBuffer(param)
	filename := fmt.Sprintf("ETH%s.json", height)

	n, err := minioClient.PutObject(bucketName, filename, buf, int64(len(param)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Info(fmt.Sprintf("upload %s size %d", filename, n))
	return nil
}

func Gzip(data []byte) []byte {
	var res bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&res, 7)
	_, err := gz.Write(data)
	if err != nil {
		fmt.Println(err)
	} else {
		gz.Close()
	}
	return res.Bytes()
}

func ParseGzip(height string) ([]byte, error) {
	data, err := GetJson(height)
	if err != nil {
		fmt.Println(err)
		return data, err
	}
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, data)
	r, err := gzip.NewReader(b)
	if err != nil {
		fmt.Println(err)
		return data, err
	} else {
		defer r.Close()
		undatas, err := ioutil.ReadAll(r)
		if err != nil {
			fmt.Println(err)
			return data, err
		}
		return undatas, nil
	}
}

func GetJson(height string) ([]byte, error) {
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	//useSSL := true
	// Initialize minio client object.
	minioClient, err := minio.New(Endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		return nil, err
	}
	// Make a new bucket called mymusic.
	bucketName := "lotus-state"
	filename := fmt.Sprintf("ETH%s.json", height)

	objectCh := minioClient.ListObjects(bucketName, filename, false, nil)
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return nil, object.Err
		}
		if object.Key == filename {
			o, err := minioClient.GetObject(bucketName, filename, minio.GetObjectOptions{})
			if err != nil {
				return nil, err
			}
			stat, err := o.Stat()
			if err != nil {
				return nil, err
			}
			data := make([]byte, stat.Size)
			_, err = o.Read(data)
			if err != nil && err != io.EOF {
				return nil, err
			}
			return data, nil
		}
	}
	return nil, errors.New("can't get height.json")
}
