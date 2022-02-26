package blockee

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type ERC20Token struct {
	//转账人
	From common.Address `json:"from"`
	//接收人
	To common.Address `json:"to"`
	//交易代币金额
	Value *big.Int `json:"value"`
	//合约地址
	ContractAddress common.Address `json:"contract"`
	//代币交易哈希
	Hash common.Hash `json:"hash"`
	//当前区块from的余额
	FromBalance *big.Int `fromBalance`
	//当前区块to的余额
	ToBalance *big.Int `toBalance`
	//代币名全称
	Name string `name`
	//代币缩写
	Symbol string `Symbol`
}
