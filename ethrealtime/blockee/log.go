package blockee

import (
	"context"
	"encoding/hex"
	"ethernum/token"
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
)

var chainId = new(big.Int).SetUint64(uint64(1))

type Log struct {
	// Consensus fields:
	// address of the contract that generated the event
	// 发生事件的合约地址
	Address common.Address `json:"address" gencodec:"required"`
	// list of topics provided by the contract.
	// 事件的主题
	Topics []common.Hash `json:"topics" gencodec:"required"`
	// supplied by the contract, usually ABI-encoded
	// 事件额外的数据
	//Data []byte `json:"data" gencodec:"required"`
	Data string `json:"data" gencodec:"required"`
	// Derived fields. These fields are filled in by the node
	// but not secured by consensus.
	// block in which the transaction was included
	// 区块编号
	BlockNumber uint64 `json:"blockNumber"`
	// hash of the transaction
	// 交易哈希
	TxHash common.Hash `json:"transactionHash" gencodec:"required"`
	// index of the transaction in the block
	// 交易索引
	TxIndex uint `json:"transactionIndex"`
	// hash of the block in which the transaction was included
	// 区块哈希
	BlockHash common.Hash `json:"blockHash"`
	// index of the log in the block
	// 日志索引
	Index uint   `json:"logIndex"`
	Time  uint64 `json:"timeStamp"`
	// The Removed field is true if this log was reverted due to a chain reorganisation.
	// You must pay attention to this field if you receive logs through a filter query.
	Removed   bool        `json:"removed"`
	ERCTokens []*ERCToken `json:"erctokens"`
	ERCType   string      `json:"ercType"`
}

func NewLog() *Log {
	log := Log{
		Address:     common.HexToAddress(""),
		Topics:      nil,
		Data:        "",
		BlockNumber: 0,
		TxHash:      common.HexToHash(""),
		TxIndex:     0,
		BlockHash:   common.HexToHash(""),
		Index:       0,
		Time:        0,
		Removed:     false,
		ERCTokens:   nil,
		ERCType:     "",
	}
	return &log
}

var Transferevent = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

const TokenABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_totalSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_TaxCollection\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"FeeWhiteList\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"TxBlackList\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"UpdateFee\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"UpdateFeeWhiteList\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"UpdateTxBlackList\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WithdrawalTRX\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"_tokenaddr\",\"type\":\"address\"}],\"name\":\"WithdrawalToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"

func (l *Log) Logs(txs types.Transactions, time uint64, client *ethclient.Client) []*Log {
	logs := make([]*Log, 0)
	contractAbi, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		fmt.Println(err)
	}
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
				//txlog := blockee.Log{}
				//erc20Tokens := make([]blockee.ERC20Token, 0)
				//日志为空，to为空，交易成功
				if (receipt.Logs == nil || len(receipt.Logs) == 0) && tx.To() == nil && receipt.Status == 1 {
					txlog := CreatorContractLog(receipt.BlockNumber, receipt.BlockHash, tx, receipt, time, client)
					logs = append(logs, txlog)
					fmt.Println()
				}

				for _, v := range receipt.Logs {

					ercTokens := make([]*ERCToken, 0)

					l.BlockHash = v.BlockHash
					//txlog.Data = v.Data
					l.Data = hex.EncodeToString(v.Data)
					l.Address = v.Address
					l.BlockNumber = v.BlockNumber
					l.Index = v.Index
					l.Topics = v.Topics
					l.TxIndex = v.TxIndex
					l.TxHash = v.TxHash
					l.Time = time
					l.Removed = v.Removed

					instance, err := token.NewToken(v.Address, client)
					if err != nil {
						log.Println("instance err:", err)
					}

					if len(v.Topics) != 3 {
						logs = append(logs, l)
						continue
					}
					ercToken := NewERCToken()
					ercToken.ERCToken(v, client, v.Topics, contractAbi, tx, instance)
					state := ercToken.TokenType(ercToken.Address, client)
					ercToken.Time = time
					ercToken.ERCType = "ERC20"
					l.ERCType = "ERC20"
					if state {
						ercToken.ERCType = "ERC721"
						l.ERCType = "ERC721"
					}
					ercTokens = append(ercTokens, ercToken)

					l.ERCTokens = ercTokens

					if len(ercToken.ContractAddress) == 0 && tx.To() == nil && receipt.Status == 1 {
						l = CreatorContractLog(receipt.BlockNumber, receipt.BlockHash, tx, receipt, time, client)
					}
					logs = append(logs, l)
				}
			}
		}(tx)
	}

	wait.Wait()
	return logs
}

func CreatorContractLog(blockNum *big.Int, blockHash common.Hash, tx *types.Transaction, receipt *types.Receipt, time uint64, client *ethclient.Client) *Log {
	//初始化
	txlog := NewLog()
	ercTokens := make([]*ERCToken, 0)

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
	instance, err := token.NewToken(txlog.Address, client)
	if err != nil {
		log.Println("instance err:", err)
		return txlog
	}

	//ercToken信息
	ercToken := NewERCToken()
	ercToken.Address = receipt.ContractAddress

	msg, err := tx.AsMessage(types.LatestSignerForChainID(chainId), tx.GasPrice())
	if err != nil {
		fmt.Printf("error %s\n", err.Error())
	}
	ercToken.ContractAddress = msg.From().String()

	state := ercToken.TokenType(ercToken.Address, client)
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
