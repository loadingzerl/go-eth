package ut_tool

import (
	"bufio"
	"encoding/json"
	"ethernum/blockee"
	"ethernum/filepath"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

var Bolckfilelist map[int64]bool
var StatusMap map[int64]bool

func MaxblockNumber(fil_path string) *big.Int {
	Bolckfilelist = make(map[int64]bool)
	//路径
	files, err := ioutil.ReadDir(fil_path)
	if err != nil {
		log.Fatal(err)
	}
	var max int64
	for _, f := range files {
		i, err := strconv.Atoi(f.Name()[:len(f.Name())-5])
		if err != nil {
			panic(err)
		}
		num := int64(i)
		if max < num {
			max = num
		}
		Bolckfilelist[num] = true
	}
	blockNumber := big.NewInt(max)
	return blockNumber
}

func MinblockNumber(fil_path string) *big.Int {
	Bolckfilelist = make(map[int64]bool)
	//路径
	files, err := ioutil.ReadDir(fil_path)
	if err != nil {
		log.Fatal(err)
	}
	i, err := strconv.Atoi(files[0].Name()[:len(files[0].Name())-5])
	if err != nil {
		panic(err)
	}
	var min = int64(i)

	for _, f := range files {
		i, err := strconv.Atoi(f.Name()[:len(f.Name())-5])
		if err != nil {
			panic(err)
		}
		num := int64(i)

		if min > num {
			min = num
		}
		Bolckfilelist[num] = true
	}
	blockNumber := big.NewInt(min)
	return blockNumber
}

func CreatedJSON(file_path string, name string, res []byte) {
	var pathFile string
	//file, err := os.Create(fmt.Sprintf("%s/%s.json", "/meta/apri/ethbtldata", name)) //创建文件
	pathFile = fmt.Sprintf("%s/%s.json", file_path, name)
	file, err := os.Create(pathFile) //创建文件
	if err != nil {
		fmt.Println("error found JSON: ", err)
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(res)
	if err != nil {
		fmt.Println("err: ---", err)
	}
}

func CreatedJSONLogs(file_path string, name string, res []byte) {
	var pathFile string
	//file, err := os.Create(fmt.Sprintf("%s/%s.json", "/meta/apri/ethbtldata", name)) //创建文件
	pathFile = fmt.Sprintf("%s/%s.json", file_path, name)
	file, err := os.Create(pathFile) //创建文件
	if err != nil {
		fmt.Println("error found JSON: ", err)
		panic(err)
	}
	defer file.Close()
	_, err = file.WriteString("[" + string(res) + "]" + "\n")
	if err != nil {
		fmt.Println("err: ---", err)
	}
}

func FileWrite(file_path string, name string, res []byte) {
	pathFile := fmt.Sprintf("%s/%s.json", file_path, name)
	file, err := os.OpenFile(pathFile, os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println("文件打开失败", err)
		//调 创建方法
		CreatedJSONLogs(file_path, name, res)
	}

	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	//write.Write(res)
	write.WriteString("[" + string(res) + "]" + "\n")
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}

func StatusMonitoring(number int64) {
	time.Sleep(3 * 60 * time.Second)
	if !StatusMap[number] {
		fmt.Println("区块:", number, " -运行超时-:")
		os.Exit(0)
	}
}

func BlockWriteFile(block *blockee.Block, txlogs []*blockee.Log, addressMap []*blockee.Address) {
	fmt.Println("blockGoroup start：", block.Number)
	resBlock, err := json.Marshal(block)
	if err != nil {
		fmt.Println("err---:", err)
	}
	n := block.Number.String()
	CreatedJSON(filepath.BlockChainPath, n, resBlock)
	fmt.Println("区块", block.Number, "写入成功")
	for _, item := range txlogs {
		resLogs, err := json.Marshal(item)
		if err != nil {
			fmt.Println("err---:", err)
		}
		FileWrite(filepath.BlockChainLogsPath, n+"_logs", resLogs)
	}
	fmt.Println("区块", block.Number, "合约监听事件写入成功")
	for _, item := range addressMap {
		addressS, err := json.Marshal(item)
		if err != nil {
			fmt.Println("err---:", err)
		}
		FileWrite(filepath.BlockChainAddressPath, n+"_address", addressS)
	}
	fmt.Println("区块", block.Number, "地址状态写入成功")
}
