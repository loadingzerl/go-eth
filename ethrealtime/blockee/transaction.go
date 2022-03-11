package blockee

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type Transaction struct {
	Hash common.Hash `json:"txHash"`
	//转账人
	From common.Address `json:"TxnFrom"`
	//接收人
	To *common.Address `json:"txnTo"`
	//执行这个交易所需要的gas。
	Gas uint64 `json:"txnGasLimt"`
	//当前gas与以太币换算的汇率。
	GasPrice *big.Int `json:"txnGasPrice"`
	//部署智能合约的交易，所以这里的input是合约的16进制代码
	Data string `json:"txnInputData"`
	//交易下的nonce值，是账户发起交易所维护的nonce，一个交易对应一个nonce值，注意区分区块中的nonce，区块中的nonce是用于POW的nonce。
	Nonce uint64 `json:"txnNonce"`
	//交易金额
	Value string `json:"txnValue"`

	Logs []Log `json:"logs"`


	ContractAddress common.Address`json:"txnContractAddress"`
	Size common.StorageSize `json:"size"`
	TxnType uint8 `json:"txnType"`
	TxnStatus uint64 `json:"txnStatus"`
	ChainId *big.Int `json:"chainId"`
	PostState  string `json:"txnPostState"`
	MaxPriority *big.Int `json:"txnMaxPriority"` //gasprice - GasFeeCap
	GasTipCap *big.Int `json:"txnGasTipCap"`
	GasFeeCap *big.Int `json:"txnGasFeeCap"`
	CumulativeGasUsed uint64 `json:"txnCumulativeGasUsed"`
	TransactionIndex uint `json:"txnTransactionIndex"`
	AccessList types.AccessList `json:"txnAccessList"`
	TxnGasUsed  uint64`json:"txGasUsed"`


}
