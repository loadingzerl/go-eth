package blockee

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Transaction struct {
	Hash common.Hash `json:"hash"`
	//转账人
	From common.Address `json:"from"`
	//接收人
	To *common.Address `json:"to"`
	//执行这个交易所需要的gas。
	Gas uint64 `json:"gas"`
	//当前gas与以太币换算的汇率。
	GasPrice *big.Int `json:"gasPrice"`
	//部署智能合约的交易，所以这里的input是合约的16进制代码
	Data []byte `json:"data"`
	//交易下的nonce值，是账户发起交易所维护的nonce，一个交易对应一个nonce值，注意区分区块中的nonce，区块中的nonce是用于POW的nonce。
	Nonce uint64 `json:"nonce"`
	//交易金额
	Value *big.Int `json:"value"`

	Logs []Log `json:"logs"`
}
