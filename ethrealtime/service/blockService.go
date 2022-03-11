package service

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"ethernum/blockee"
	"ethernum/filepath"
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
	"strings"
	"sync"
	"time"
)

var Client *ethclient.Client

var Transferevent = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

const TokenABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_totalSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_TaxCollection\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"FeeWhiteList\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"TxBlackList\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"UpdateFee\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"UpdateFeeWhiteList\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"UpdateTxBlackList\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WithdrawalTRX\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"_tokenaddr\",\"type\":\"address\"}],\"name\":\"WithdrawalToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"


func ClientInit() {
	for i := 0; i <= 9; i++ {
		var err error
		//client, err = ethclient.Dial("wss://mainnet.infura.io/ws/v3/e2469954ed8248698d0c6d10d1127dd2")
		//韩国
		//Client, err = ethclient.Dial("ws://138.113.235.159:8546")
		//Client, err = ethclient.Dial("ws://138.113.235.159:8546")

		//深圳
		Client, err = ethclient.Dial("ws://183.60.141.2:8546")
		if err != nil {
			time.Sleep(5 * time.Second)
			log.Println("client_err:", err)
			continue
		}
		break
	}

}


func HistoryBlock(minNmuber *big.Int,headerNumber *big.Int)  {
	header, err := Client.HeaderByNumber(context.Background(), nil)
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

func BlockGoroup(number *big.Int) {
	start := time.Now() // 获取当前时间
	block, err := Client.BlockByNumber(context.Background(), number)
	fmt.Println("blockGoroup start：", block.Number())
	blockMap, logs := BlockNumber(block)

	addressS := AddressStatus(block)
	//写入文件
	resBlock, err := json.Marshal(blockMap)
	if err != nil {
		fmt.Println("err---:", err)
	}
	//resLogs, err := json.Marshal(logs)
	//if err != nil {
	//	fmt.Println("err---:", err)
	//}
	n := blockMap.Number.String()
	//fmt.Println(n)
	//createdJSON(n, resBlock, 1)

	util.CreatedJSON(filepath.BlockChainPath, n, resBlock)

	fmt.Println("区块", blockMap.Number, "写入成功")
	//createdJSON(n+"_logs", resLogs, 0)

	for _, item := range logs {
		resLogs, err := json.Marshal(item)
		//_, err := json.Marshal(item)
		if err != nil {
			fmt.Println("err---:", err)
		}
		util.FileWrite(filepath.BlockChainLogsPath, n+"_logs", resLogs)
	}

	for _, item := range addressS {
		addressS, err := json.Marshal(item)
		//_, err := json.Marshal(item)

		if err != nil {
			fmt.Println("err---:", err)
		}
		util.FileWrite(filepath.BlockChainAddressPath, n+"_address", addressS)
	}
	//util.CreatedJSON(BlockChainLogsPath, n+"_logs", resLogs)
	fmt.Println("区块", blockMap.Number, "合约监听事件写入成功")




	fmt.Println("blockGoroup end: ", blockMap.Number)
	fmt.Println("end : ", time.Now().Format("2006-01-02 15:04:05"))
	elapsed := time.Since(start)
	fmt.Println("该函数执行完成耗时：", elapsed)

}

func BlockNumber(block *types.Block) (blockee.Block, []blockee.Log) {
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
	blockMap.Extra = hex.EncodeToString(block.Extra())
	blockMap.Hash = block.Hash()
	blockMap.MixDigest = block.MixDigest()
	blockMap.Nonce = block.Nonce()
	blockMap.Size = block.Size()
	blockMap.ReceivedFrom = block.ReceivedFrom
	blockMap.UnclesNumber = len(block.Uncles())
	blockMap.ReceivedAa = block.ReceivedAt
	blockMap.Coinbase = block.Coinbase()
	blockMap.Uncles = block.Uncles()
	blockMap.BlockBaseFee = block.BaseFee()

	//遍历交易信息t
	transactions := Transactions(block)
	blockMap.Transactions = transactions
	txs := block.Transactions()
	logs := Logs(txs,block.Time())
	//fmt.Println(logs,"KKKKKKKKKKKKKKKKKKKKKK")
	return blockMap, logs
}
var chainId = new(big.Int).SetUint64(uint64(1))
func AddressStatus(block *types.Block) []blockee.Address  {
	addressStatus := make([]blockee.Address,0)
	for _,tx := range block.Transactions() {
		addressFrom := blockee.Address{}
		addressTo := blockee.Address{}
		msg, err := tx.AsMessage(types.LatestSignerForChainID(chainId),tx.GasPrice())
		if err == nil {
			addressFrom.Address = msg.From()
		}
		balanceFrom ,_:= Client.BalanceAt(context.Background(),addressFrom.Address,block.Number())
		addressFrom.Balance = balanceFrom.String()
		contractAbi, err := abi.JSON(strings.NewReader(TokenABI))
		if err != nil {
			fmt.Println(err)
		}
		addressFrom.Code = contractAbi
		addressFrom.Nonce = tx.Nonce()
		addressFrom.Number =block.Number()
		addressFrom.Time =block.Time()

		stotagekeyFrom, err:= Client.StorageAt(context.Background(),addressFrom.Address,tx.Hash(),block.Number())
		if err != nil {
			fmt.Println(err)
		}
		addressFrom.Storagekey = hex.EncodeToString(stotagekeyFrom)
		addressStatus = append(addressStatus,addressFrom)

		if tx.To() != nil {
			addressTo.Address = *tx.To()
		}
		balanceTo ,_:= Client.BalanceAt(context.Background(),addressFrom.Address,block.Number())
		addressTo.Balance = balanceTo.String()
		addressTo.Code = contractAbi
		addressTo.Nonce = tx.Nonce()
		addressTo.Number = block.Number()
		addressTo.Time = block.Time()
		stotagekeyTo, err:= Client.StorageAt(context.Background(),addressTo.Address,tx.Hash(),block.Number())
		if err != nil {
			fmt.Println(err)
		}
		addressTo.Storagekey = hex.EncodeToString(stotagekeyTo)
		addressStatus = append(addressStatus,addressTo)
	}
	return addressStatus
}

func Transactions(block *types.Block) []blockee.Transaction {
	Transactions := make([]blockee.Transaction, 0)
	var chainId = new(big.Int).SetUint64(uint64(1))
	for _, tx := range block.Transactions() {
		transaction := blockee.Transaction{}
		transaction.Hash = tx.Hash()
		msg, err := tx.AsMessage(types.LatestSignerForChainID(chainId),tx.GasPrice())
		if err == nil {
			transaction.From = msg.From()
		}
		transaction.To = tx.To()
		if transaction.To == nil {
			to := crypto.CreateAddress(msg.From(),msg.Nonce())
			transaction.ContractAddress = to
		}
		transaction.Gas = tx.Gas()
		transaction.Data = hex.EncodeToString(tx.Data())
		transaction.Nonce = tx.Nonce()
		transaction.Value = tx.Value().String()
		receipt, err := Client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			fmt.Println("tx receipt ",receipt)
		}
		if receipt != nil {
			transaction.ContractAddress = receipt.ContractAddress
			transaction.TxnType = receipt.Type
			transaction.TxnStatus = receipt.Status
			transaction.PostState = hex.EncodeToString(receipt.PostState)
			transaction.TxnGasUsed = receipt.GasUsed
			transaction.CumulativeGasUsed = receipt.CumulativeGasUsed
			transaction.TransactionIndex = receipt.TransactionIndex
			transaction.PostState = hex.EncodeToString(receipt.PostState)
		}
		transaction.Size = tx.Size()
		transaction.ChainId = tx.ChainId()
		txgas ,_,err :=Client.TransactionByHash(context.Background(),tx.Hash())
		if err != nil{
			fmt.Println(err)
		}
		transaction.GasPrice = msg.GasPrice()
		gas := big.Int{}
		transaction.GasFeeCap = gas.Sub(transaction.GasPrice,txgas.GasTipCap())
		transaction.MaxPriority = tx.GasTipCap()
		transaction.GasTipCap = tx.GasTipCap()

		//if transaction.TxnType==2{
		//	fmt.Println("type :",transaction.TxnType)
		//	fmt.Println("txhash: ",tx.Hash())
		//	fmt.Println("transaction.GasPrice :", transaction.GasPrice)
		//	fmt.Println("transaction.GasFeeCap :", transaction.GasFeeCap ,"GAS use:",tx.Gas())
		//	fmt.Println("MaxPriority :",transaction.MaxPriority)
		//	fmt.Println("basefee",block.Header().BaseFee)
		//	//os.Exit(1)
		//}


		transaction.AccessList = tx.AccessList()

		Transactions = append(Transactions, transaction)
	}
	//fmt.Println("13907335-transactions: ",len(Transactions))
	return Transactions
}

func Logs(txs types.Transactions,time uint64) []blockee.Log {
	logs := make([]blockee.Log, 0)
	//erc20Tokens := make([]blockee.ERC20Token, 0)
	//var chainId = new(big.Int).SetUint64(uint64(1))
	contractAbi, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		fmt.Println(err)
	}

	//transaction := blockee.Transaction{}
	wait := sync.WaitGroup{}
	wait.Add(len(txs))
	for _, tx := range txs {
		//if tx.Hash()==common.HexToHash("0x7033af97111a5ae2673196fed2735308cadc360deac7d3dc323fa6f5d82c71fa"){
		//wait.Add(1)
		//fmt.Println("0x7033af97111a5ae2673196fed2735308cadc360deac7d3dc323fa6f5d82c71fa",i)
		go func(tx *types.Transaction) {
			defer wait.Done()
			//receiptInit:
			receipt, err := Client.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil {
				fmt.Println("receipt--err: ", err)
				fmt.Println("receipt-tx.Hash: ", tx.Hash())
				fmt.Println("receipt-->", receipt)
				//clientInit()
				//goto receiptInit
			}

			if receipt != nil {
				//txlog := blockee.Log{}
				//erc20Tokens := make([]blockee.ERC20Token, 0)
				//日志为空，to为空，交易成功
				if (receipt.Logs == nil || len(receipt.Logs) == 0) && tx.To() == nil && receipt.Status == 1 {
					fmt.Println(" receipt.Logs == nil && tx.To() == nil && receipt.Status == 1 ---------------")
					txlog := CreatorContractLog(receipt.BlockNumber, receipt.BlockHash, tx, receipt, time)
					logs = append(logs, txlog)
					fmt.Println()
				}

				for _, v := range receipt.Logs {
					ercTokens := make([]blockee.ERCToken, 0)
					txlog := blockee.Log{}
					txlog.BlockHash = v.BlockHash
					//txlog.Data = v.Data
					txlog.Data = hex.EncodeToString(v.Data)
					txlog.Address = v.Address
					txlog.BlockNumber = v.BlockNumber
					txlog.Index = v.Index
					txlog.Topics = v.Topics
					txlog.TxIndex = v.TxIndex
					txlog.TxHash = v.TxHash
					txlog.Time = time
					txlog.Removed = v.Removed
					instance, err := token.NewToken(v.Address, Client)
					if err != nil {
						log.Println("instance err:", err)
					}

					//if len(v.Topics) != 3 {
					//	logs = append(logs, txlog)
					//	continue
					//}
					ercToken := blockee.ERCToken{}

					switch v.Topics[0] {

					case Transferevent:
						event, err := contractAbi.Unpack("Transfer", v.Data)
						if err != nil {
							//log.Println("event err:", err)
							fmt.Println("event err:", err)
						}
						//ercToken := blockee.ERCToken{}
						if len(event) > 0 {
							//	"value:", event[0], "contractAddr:", v.Address, "TXhash:", tx.Hash())
							var vl *big.Int
							if value, ok := event[0].(*big.Int); ok {
								vl = value
							}
							if len(v.Topics) >=3 {
								ercToken.From = common.HexToAddress(v.Topics[1].Hex())
								ercToken.To = common.HexToAddress(v.Topics[2].Hex())
							}


							//ercToken.ContractAddress = receipt.ContractAddress.String()
							//if receipt.ContractAddress == common.HexToAddress("0x0000000000000000000000000000000000000000"){
							//	ercToken.ContractAddress = ""
							//}

							if tx.To() == nil {
								msg, err := tx.AsMessage(types.LatestSignerForChainID(chainId), tx.GasPrice())
								if err != nil {
									fmt.Printf("error %s\n", err.Error())
								}
								//合约创建地址
								ercToken.ContractAddress = msg.From().String()
							}

							ercToken.Value = vl.String()
							ercToken.Address = v.Address
							ercToken.Hash = tx.Hash()
						} else {
							if len(v.Topics) >=3 {
								ercToken.From = common.HexToAddress(v.Topics[1].Hex())
								ercToken.To = common.HexToAddress(v.Topics[2].Hex())
							}

							ercToken.Address = v.Address

							//ercToken.ContractAddress = receipt.ContractAddress.String()
							//if receipt.ContractAddress == common.HexToAddress("0x0000000000000000000000000000000000000000"){
							//	ercToken.ContractAddress = ""
							//}

							if tx.To() == nil {
								msg, err := tx.AsMessage(types.LatestSignerForChainID(chainId), tx.GasPrice())
								if err != nil {
									fmt.Printf("error %s\n", err.Error())
								}
								//合约创建地址
								ercToken.ContractAddress = msg.From().String()
							}

							ercToken.Hash = tx.Hash()
						}
						name, err := instance.Name(nil)
						if err != nil {
							fmt.Println("ERCTOKENname err:", err)
						}
						Symbol, err := instance.Symbol(nil)
						if err != nil {
							fmt.Println("ERCTOKENSymbol err:", err)
						}
						addressBalance ,err := Client.BalanceAt(context.Background(),ercToken.Address,nil)
						if err != nil {
							fmt.Println(err)
						}
						ercToken.AddressBalance = addressBalance.String()
						var fromBalance *big.Int
						var toBalance *big.Int
						if len(v.Topics) >=3 {
							fromBalance, err = instance.BalanceOf(nil, common.HexToAddress(v.Topics[1].Hex()))
							if err != nil {
								fmt.Println("fromBalance", err)
							}
							toBalance, err = instance.BalanceOf(nil, common.HexToAddress(v.Topics[2].Hex()))
							if err != nil {
								fmt.Println("toBalance", err)
							}
						}

						ercToken.BlockNumber = receipt.BlockNumber
						ercToken.BlockHash = receipt.BlockHash
						decimals, err  := instance.Decimals(nil)
						if err != nil {
							fmt.Println("decimals err ", err)
						}
						ercToken.Time = time
						ercToken.Decimals = decimals
						ercToken.FromBalance = fromBalance.String()
						ercToken.ToBalance = toBalance.String()
						ercToken.Name = name
						ercToken.Symbol = Symbol
						state := TokenType(ercToken.Address)
						ercToken.ERCType = "ERC20"
						txlog.ERCType = "ERC20"
						if state {
							ercToken.ERCType = "ERC721"
							txlog.ERCType = "ERC721"
						}
						ercTokens = append(ercTokens, ercToken)

						txlog.ERCTokens = ercTokens
					}

					if  len(ercToken.ContractAddress) == 0 && tx.To() == nil && receipt.Status == 1 {
						txlog = CreatorContractLog(receipt.BlockNumber, receipt.BlockHash, tx, receipt, time)
					}

					logs = append(logs, txlog)
				}

			}
		}(tx)
	}

	wait.Wait()
	return logs
}

func TokenType(contractAddress common.Address) bool  {
	erc721contract,err:=token.NewErc721(contractAddress,Client)
	if err != nil {
		fmt.Println(err)
	}
	data := [4]byte{0x80,0xac,0x58,0xcd}
	iserc721,err:=erc721contract.SupportsInterface(nil,data)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("查询ERC721结果:",iserc721)
	return iserc721
}

func CreatorContractLog(blockNum *big.Int, blockHash common.Hash, tx *types.Transaction, receipt *types.Receipt, time uint64) blockee.Log {
	//初始化
	txlog := blockee.Log{}
	ercTokens := make([]blockee.ERCToken, 0)

	//txlog赋值
	txlog.Address = receipt.ContractAddress
	//txlog.Topics
	//txlog.Data
	txlog.TxHash = tx.Hash()
	txlog.TxIndex = receipt.TransactionIndex
	txlog.BlockHash = blockHash
	txlog.BlockNumber = blockNum.Uint64()
	txlog.Index = 4380000
	//txlog.ERCType
	//txlog.Removed
	txlog.Time = time

	//erc赋值
	instance, err := token.NewToken(txlog.Address, Client)
	if err != nil {
		log.Println("instance err:", err)
		return txlog
	}

	//ercToken信息
	ercToken := blockee.ERCToken{}
	ercToken.Address = receipt.ContractAddress

	msg, err := tx.AsMessage(types.LatestSignerForChainID(chainId), tx.GasPrice())
	if err != nil {
		fmt.Printf("error %s\n", err.Error())
	}
	ercToken.ContractAddress = msg.From().String()

	state := TokenType(ercToken.Address)
	if state {
		ercToken.ERCType = "ERC721"
		txlog.ERCType = "ERC721"
	} else {
		ercToken.ERCType = "ERC20"
		txlog.ERCType = "ERC20"
	}
	ercToken.BlockNumber = blockNum
	ercToken.BlockHash = blockHash

	name, err := instance.Name(nil)
	if err != nil {
		fmt.Println("ERCTOKENname err:", err)
	}
	symbol, err := instance.Symbol(nil)
	if err != nil {
		fmt.Println("ERCTOKENSymbol err:", err)
	}
	ercToken.Name = name
	ercToken.Symbol = symbol

	decimals, err := instance.Decimals(nil)
	if err != nil {
		fmt.Println("decimals err ", err)
	}
	ercToken.Decimals = decimals
	ercToken.Time = time

	ercTokens = append(ercTokens, ercToken)

	txlog.ERCTokens = ercTokens

	return txlog
}