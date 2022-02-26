package main

import (
	"context"
	"encoding/json"
	"ethernum/blockee"
	"ethernum/token"
	"ethernum/util"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"
)

var client *ethclient.Client
var statusMap map[int64]bool

const (
	//BlockChainPath     string = "/home/loading/goProject/src/ethrealtime/blockChain"
	//BlockChainLogsPath string = "/home/loading/goProject/src/ethrealtime/blockChainLogs"
	//BlockChainDataPath string = "/home/loading/goProject/src/ethrealtime/blockChain/"

	//BlockChainPath     string = "/meta/apri/ethblockdata"
	//BlockChainLogsPath string = "/meta/apri/etherc20blockdata"
	//BlockChainDataPath string = "/meta/apri/ethblockdata/"

	BlockChainPath     string = "/meta/ethrealtime/ethblockdata"
	BlockChainLogsPath string = "/meta/ethrealtime/etherc20blockdata"
	BlockChainDataPath string = "/meta/ethrealtime/ethblockdata/"
)

type EthController struct {
}

var Transferevent = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

const TokenABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_totalSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_TaxCollection\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"FeeWhiteList\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"TxBlackList\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"UpdateFee\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"UpdateFeeWhiteList\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"UpdateTxBlackList\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WithdrawalTRX\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"_tokenaddr\",\"type\":\"address\"}],\"name\":\"WithdrawalToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"

func clientInit() {
	for i := 0; i <= 9; i++ {
		var err error
		//client, err = ethclient.Dial("wss://mainnet.infura.io/ws/v3/e2469954ed8248698d0c6d10d1127dd2")
		//韩国
		//client, err = ethclient.Dial("ws://138.113.235.159:8546")
		//深圳 ws://183.60.141.2:8546
		client, err = ethclient.Dial("ws://138.113.235.159:854")
		if err != nil {
			time.Sleep(5 * time.Second)
			log.Println("client_err:", err)
			continue
		}
		break
	}

}

func historyBlock(minNmuber *big.Int,headerNumber *big.Int)  {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(header.Number)
	endNumber:=headerNumber
	minNumber := minNmuber.Int64()
	//minNumber := int64(8618924)
	c := make(chan int64, 10)
	wait := sync.WaitGroup{}
	//wait.Add(1)
	fmt.Println("-------:",endNumber.Int64()-int64(minNumber)-1)
	for i:=minNumber;i<endNumber.Int64();i++{
		if !util.Bolckfilelist[i]{
			//fmt.Println("----",i)
			c<-i
			go func(blockNmuber int64) {
				wait.Add(1)
				blockGoroup(big.NewInt(blockNmuber))
				//time.Sleep(time.Second)
				<-c
				wait.Done()
				//fmt.Println("写入成功",blockNmuber)
			}(i)
		}
	}

	wait.Wait()
}

func main() {
	clientInit()

	header, err := client.HeaderByNumber(context.Background(), nil)
	var munNumber = util.MinblockNumber(BlockChainDataPath)
	if munNumber.Int64() < header.Number.Int64() {
		historyBlock(munNumber,header.Number)
	}

	fmt.Println("---缺少的历史区块补齐成功---")
    statusMap = make(map[int64]bool)
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	fmt.Println("sub = ", sub)
	if err != nil {
		log.Println("ErrorSub:", err)
		os.Exit(0)
	}

	//blockGoroup(big.NewInt(13936616))

	for {
		select {
		case err := <-sub.Err():
			log.Println("sub->:", err)
		case header := <-headers:
			statusMap[header.Number.Int64()-2] =  false
			go statusMonitoring(header.Number.Int64()-2)
			//blockGoroup+(big.NewInt(number))
			blockGoroup(big.NewInt(header.Number.Int64()-2))
			fmt.Println("区块： ",header.Number.Int64()-2 ,"：----Execution complete----")
			statusMap[header.Number.Int64()-2] = true
		}
	}

}

func statusMonitoring(number int64)  {
	time.Sleep(3*60*time.Second)
	if !statusMap[number] {
		fmt.Println("区块:",number," -运行超时-:")
		os.Exit(0)
	}
}

func blockGoroup(number *big.Int) {

	start := time.Now() // 获取当前时间
	eth := EthController{}
	block, err := client.BlockByNumber(context.Background(), number)
	fmt.Println("blockGoroup start：", block.Number())
	blockMap, logs := eth.BlockNumber(block)
	//写入文件
	resBlock, err := json.Marshal(blockMap)
	if err != nil {
		fmt.Println("err---:", err)
	}

	n := blockMap.Number.String()
	//createdJSON(n, resBlock, 1)
	util.CreatedJSON(BlockChainPath, n, resBlock)
	fmt.Println("区块", blockMap.Number, "写入成功")
	//fmt.Println("13907335->block-transactions: ",len(blockMap.Transactions))
	//createdJSON(n+"_logs", resLogs, 0)
	for _, item := range logs {
		resLogs, err := json.Marshal(item)
		if err != nil {
			fmt.Println("err---:", err)
		}
		util.FileWrite(BlockChainLogsPath, n+"_logs", resLogs)
	}
	//util.CreatedJSON(BlockChainLogsPath, n+"_logs", resLogs)
	fmt.Println("区块", blockMap.Number, "合约监听事件写入成功")

	fmt.Println("blockGoroup end: ", blockMap.Number)
	fmt.Println("end : ", time.Now().Format("2006-01-02 15:04:05"))
	elapsed := time.Since(start)
	fmt.Println("该函数执行完成耗时：", elapsed)


}

func (eth EthController) BlockNumber(block *types.Block) (blockee.Block, []blockee.Log) {
//blockInit:
//	block, err := client.BlockByNumber(context.Background(), number)
//	if err != nil {
//		fmt.Println("BlockNumber() err : ", err)
//		clientInit()
//		//block, err = client.BlockByNumber(context.Background(), number)
//		//goto blockInit
//	}
	//创建block对象
	blockMap := blockee.Block{}
	blockMap.Number = block.Number()
	blockMap.Time = block.Time()
	blockMap.ParentHash = block.ParentHash()
	blockMap.UncleHash = block.UncleHash()
	blockMap.Root = block.Root()
	blockMap.TxHash = block.TxHash()
	blockMap.ReceiptHash = block.ReceiptHash()
	blockMap.Difficulty = block.Difficulty()
	blockMap.GasLimit = block.GasLimit()
	blockMap.GasUsed = block.GasUsed()
	blockMap.Extra = block.Extra()
	blockMap.Hash = block.Hash()
	blockMap.MixDigest = block.MixDigest()
	blockMap.Nonce = block.Nonce()
	blockMap.Size = block.Size()


	//遍历交易信息t
	transactions := eth.Transactions(block)
	blockMap.Transactions = transactions
	txs := block.Transactions()
	logs := eth.Logs(txs,block.Time())

	return blockMap, logs
}

func (s EthController) Transactions(block *types.Block) []blockee.Transaction {
	Transactions := make([]blockee.Transaction, 0)
	var chainId = new(big.Int).SetUint64(uint64(1))
	transaction := blockee.Transaction{}
	for _, tx := range block.Transactions() {
		transaction.Hash = tx.Hash()
		msg, err := tx.AsMessage(types.LatestSignerForChainID(chainId),tx.GasPrice())
		if err == nil {
			transaction.From = msg.From()
		}
		transaction.To = tx.To()
		if transaction.To == nil {
			to := crypto.CreateAddress(msg.From(),msg.Nonce())
			transaction.To = &to
		}
		transaction.Gas = tx.Gas()
		transaction.GasPrice = tx.GasPrice()
		transaction.Data = tx.Data()
		transaction.Nonce = tx.Nonce()
		transaction.Value = tx.Value()
		Transactions = append(Transactions, transaction)
	}
	//fmt.Println("13907335-transactions: ",len(Transactions))
	return Transactions
}

func (s EthController) Logs(txs types.Transactions,time uint64) []blockee.Log {
	logs := make([]blockee.Log, 0)
	erc20Tokens := make([]blockee.ERC20Token, 0)
	//var chainId = new(big.Int).SetUint64(uint64(1))
	contractAbi, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		fmt.Println(err)
	}
	//transaction := blockee.Transaction{}
	wait := sync.WaitGroup{}
	wait.Add(len(txs))
	for _, tx := range txs {
		go func(tx *types.Transaction) {
			defer wait.Done()
			//receiptInit:
			receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil {
				fmt.Println("receipt--err: ", err)
				fmt.Println("receipt-tx.Hash: ", tx.Hash())
				fmt.Println("receipt-->", receipt)
				//clientInit()
				//goto receiptInit
			}
			if receipt != nil {
				for _, v := range receipt.Logs {
					txlog := blockee.Log{}
					txlog.BlockHash = v.BlockHash
					txlog.Data = v.Data
					txlog.Address = v.Address
					txlog.BlockNumber = v.BlockNumber
					txlog.Index = v.Index
					txlog.Topics = v.Topics
					txlog.TxIndex = v.TxIndex
					txlog.TxHash = v.TxHash
					txlog.Time = time
					instance, err := token.NewToken(v.Address, client)
					if err != nil {
						log.Println("instance err:", err)
						fmt.Println("instance err:", err)
					}

					if len(v.Topics) < 3 {
						continue
					}
					switch v.Topics[0] {

					case Transferevent:
						event, err := contractAbi.Unpack("Transfer", v.Data)
						if err != nil {
							//log.Println("event err:", err)
							fmt.Println("event err:", err)
						}
						erc20Token := blockee.ERC20Token{}
						if len(event) > 0 {
							//	"value:", event[0], "contractAddr:", v.Address, "TXhash:", tx.Hash())
							var vl *big.Int
							if value, ok := event[0].(*big.Int); ok {
								vl = value
							}
							//if len(v.Topics) >= 3 {
							erc20Token.From = common.HexToAddress(v.Topics[1].Hex())
							erc20Token.To = common.HexToAddress(v.Topics[2].Hex())
							//}
							//erc20Token.From = common.HexToAddress(v.Topics[1].Hex())
							//erc20Token.To = common.HexToAddress(v.Topics[2].Hex())
							erc20Token.Value = vl
							erc20Token.ContractAddress = v.Address
							erc20Token.Hash = tx.Hash()
							//fmt.Println("value:", event[0])
						} else {
							//if len(v.Topics) >= 3 {
							erc20Token.From = common.HexToAddress(v.Topics[1].Hex())
							erc20Token.To = common.HexToAddress(v.Topics[2].Hex())
							//}
							erc20Token.ContractAddress = v.Address
							erc20Token.Hash = tx.Hash()
						}
						name, err := instance.Name(nil)
						if err != nil {
							//log.Println("ERCTOKENname err:", err)
							fmt.Println("ERCTOKENname err:", err)
						}
						Symbol, err := instance.Symbol(nil)
						if err != nil {
							//log.Println("ERCTOKENSymbol err:", err)
							fmt.Println("ERCTOKENSymbol err:", err)
						}

						var fromBalance *big.Int
						var toBalance *big.Int
						//if len(v.Topics) >= 3 {
						fromBalance, err = instance.BalanceOf(nil, common.HexToAddress(v.Topics[1].Hex()))
						if err != nil {
							fmt.Println("fromBalance", err)
						}
						toBalance, err = instance.BalanceOf(nil, common.HexToAddress(v.Topics[2].Hex()))
						if err != nil {
							fmt.Println("toBalance", err)
						}
						//}
						erc20Token.FromBalance = fromBalance
						erc20Token.ToBalance = toBalance
						erc20Token.Name = name
						erc20Token.Symbol = Symbol
						erc20Tokens = append(erc20Tokens, erc20Token)
						txlog.ERC20Tokens = erc20Tokens
					}

					logs = append(logs, txlog)
				}
			}
		}(tx)
	}
	wait.Wait()
	return logs
}
