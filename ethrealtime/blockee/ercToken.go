package blockee

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type ERCToken struct {
	//转账人
	From common.Address `json:"from"`
	//接收人
	To common.Address `json:"to"`
	//交易代币金额
	Value string `json:"value"`
	//创建合约的地址 from
	ContractAddress string `json:"creator"`
	//代币交易哈希
	Hash common.Hash `json:"hash"`
	//当前区块from的余额
	FromBalance string `json:"fromBalance"`
	//当前区块to的余额
	ToBalance string `json:"toBalance"`
	//代币名全称
	Name string `json:"name"`
	//代币缩写
	Symbol string `json:"symbol"`

	//合约地址
	Address  common.Address `json:"address"`
	BlockHash common.Hash `json:"blockHash"`
	BlockNumber *big.Int `json:"blockNumber"`
	Decimals uint8 `json:"decimals"`
	ERCType string `json:"ercType"`
	Time uint64 `json:"timeStamp"`
	AddressBalance string `json:"addressBalance"`
}
