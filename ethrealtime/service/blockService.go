package service

import (
	"context"
	"ethernum/filepath"
	"ethernum/modul"
	ut_tool "ethernum/tool"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"os"
	"sync"
	"time"
)

var client *ethclient.Client

func ClientInit() {
	for i := 0; i <= 9; i++ {
		var err error
		//client, err = ethclient.Dial("wss://mainnet.infura.io/ws/v3/e2469954ed8248698d0c6d10d1127dd2")

		//Client, err = ethclient.Dial("ws://138.113.235.159:8546")

		client, err = ethclient.Dial("ws://183.60.141.2:8546")
		if err != nil {
			time.Sleep(5 * time.Second)
			log.Println("client_err:", err)
			continue
		}
		break
	}

}

//realtimeBlockmonitor
func BlockTime() {
	ClientInit()
	ut_tool.StatusMap = make(map[int64]bool)
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	fmt.Println("sub = ", sub)
	if err != nil {
		log.Println("ErrorSub:", err)
		os.Exit(0)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Println("sub->:", err)
		case header := <-headers:
			ut_tool.StatusMap[header.Number.Int64()-2] = false
			go ut_tool.StatusMonitoring(header.Number.Int64() - 2)
			//blockGoroup+(big.NewInt(number))
			BlockGoroup(big.NewInt(header.Number.Int64() - 2))
			fmt.Println("区块： ", header.Number.Int64()-2, "：----Execution complete----")
			ut_tool.StatusMap[header.Number.Int64()-2] = true
		}
	}
}

func BlockGoroup(number *big.Int) {
	block1, err := client.BlockByNumber(context.Background(), number)
	if err != nil {
		fmt.Println(err)
	}
	block := modul.NewBlock()
	block.BlockObtain(block1, client)

	txlog := modul.NewLog()
	txlogs := txlog.Logs(block1.Transactions(), block1.Time(), client)

	address := modul.NewAddress()
	addressMap := address.AddressStatus(block1, client)

	start := time.Now() // 获取当前时间
	fmt.Println("blockGoroup start：", block.Number)
	//写入操作
	ut_tool.BlockWriteFile(block, txlogs, addressMap)
	fmt.Println("blockGoroup end: ", block.Number)
	fmt.Println("end : ", time.Now().Format("2006-01-02 15:04:05"))
	elapsed := time.Since(start)
	fmt.Println("该函数执行完成耗时：", elapsed)
}

//区块丢失校验
func BlockLoseCheck() {
	header, err := client.HeaderByNumber(context.Background(), nil)
	fmt.Println(header.Number)
	if err != nil {
		log.Println(err)
	}
	endNumber := int64(1400000)
	var munNumber = ut_tool.MinblockNumber(filepath.BlockChainDataPath)
	//header.Number
	if munNumber.Int64() < endNumber {
		fmt.Println("------------------------")
		HistoryBlock(munNumber, big.NewInt(endNumber))
	}
	fmt.Println("ok--------")
	fmt.Println("---缺少的历史区块补齐成功---")
}

func HistoryBlock(minNmuber *big.Int, headerNumber *big.Int) {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(header.Number)
	endNumber := headerNumber
	minNumber := minNmuber.Int64()
	//minNumber := int64(8618924)
	c := make(chan int64, 10)
	wait := sync.WaitGroup{}
	//wait.Add(1)
	fmt.Println("-------:", endNumber.Int64()-int64(minNumber)-1)
	for i := minNumber; i < endNumber.Int64(); i++ {
		if !ut_tool.Bolckfilelist[i] {
			//fmt.Println("----",i)
			c <- i
			go func(blockNmuber int64) {
				wait.Add(1)
				BlockGoroup(big.NewInt(blockNmuber))
				//time.Sleep(time.Second)
				<-c
				wait.Done()
				//fmt.Println("写入成功",blockNmuber)
			}(i)
		}
	}

	wait.Wait()
}
