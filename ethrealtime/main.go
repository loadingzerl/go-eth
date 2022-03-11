package main

import (
	"context"
	"ethernum/service"
	"ethernum/util"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math/big"
	"os"
)

func main() {
	service.ClientInit()
	client := service.Client
	//header, err := client.HeaderByNumber(context.Background(), nil)
	//var munNumber = util.MinblockNumber(filepath.BlockChainDataPath)
	//if munNumber.Int64() < header.Number.Int64() {
	//	service.HistoryBlock(munNumber,header.Number)
	//}
	//
	//fmt.Println("---缺少的历史区块补齐成功---")


	util.StatusMap = make(map[int64]bool)
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	fmt.Println("sub = ", sub)
	if err != nil {
		log.Println("ErrorSub:", err)
		os.Exit(0)
	}

	//service.BlockGoroup(big.NewInt(14356527 ))

	for {
		select {
		case err := <-sub.Err():
			log.Println("sub->:", err)
		case header := <-headers:
			util.StatusMap[header.Number.Int64()-2] =  false
			go util.StatusMonitoring(header.Number.Int64()-2)
			//blockGoroup+(big.NewInt(number))
			service.BlockGoroup(big.NewInt(header.Number.Int64()-2))
			fmt.Println("区块： ",header.Number.Int64()-2 ,"：----Execution complete----")
			util.StatusMap[header.Number.Int64()-2] = true
		}
	}

}


