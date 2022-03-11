package blockee

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"time"
)

type Block struct {

	//区块高度
	Number *big.Int `json:"blockHeiget"           gencodec:"required"`
	//时间戳
	//Time
	Time uint64 `json:"timeStamp"        gencodec:"required"`
	//交易
	Transactions []Transaction
	//挖矿人

	//前一个区块的哈希值。
	ParentHash common.Hash `json:"parentHash"       gencodec:"required"`
	//数据块的哈希值
	UncleHash common.Hash `json:"sha3Uncles"       gencodec:"required"`
	//区块状态树的根哈希
	Root common.Hash `json:"stateRoot"        gencodec:"required"`
	//区块的交易树的根哈希
	TxHash common.Hash `json:"transactionsRoot" gencodec:"required"`
	//收据树的根哈希值
	ReceiptHash common.Hash `json:"receiptsRoot"     gencodec:"required"`
	//当前区块的难度。
	Difficulty *big.Int `json:"difficulty"       gencodec:"required"`
	//当前区块允许使用的最大gas。
	GasLimit uint64 `json:"gasLimit"         gencodec:"required"`
	//当前区块累计使用的gas。
	GasUsed uint64 `json:"BlockGasUsed"          gencodec:"required"`
	//与此区块相关的附加数据。
	Extra string `json:"extraData"        gencodec:"required"`
	//区块哈希
	Hash common.Hash `json:"blockHash"`
	//一个Hash值，当与nonce组合时，证明此区块已经执行了足够的计算。 nonce：POW生成的哈希值。
	//MixDigest
	MixDigest common.Hash `json:"mixDisgest"`
	//Nonce
	Nonce     uint64      `json:"blockNonce"`
	//Size
	Size common.StorageSize `json:"blockSize"`


	ReceivedFrom  interface{} `json:"receivedFrom"`
	UnclesNumber  int `json:"unclesNumber"`
	ReceivedAa  time.Time `json:"receivedAt"`
	Coinbase  common.Address `json:"miner"            gencodec:"required"`
	Uncles []*types.Header `json:"uncles"`
	BlockBaseFee *big.Int `json:"blockBaseFee"`
}





