package blockee

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Block struct {

	//区块高度
	Number *big.Int `json:"number"           gencodec:"required"`
	//时间戳
	Time uint64 `json:"timestamp"        gencodec:"required"`
	//交易
	Transactions []Transaction
	//挖矿人

	//前一个区块的哈希值。
	ParentHash common.Hash `json:"parentHash"       gencodec:"required"`
	//数据块的哈希值
	UncleHash common.Hash `json:"sha3Uncles"       gencodec:"required"`
	//Coinbase  common.Address `json:"miner"            gencodec:"required"`
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
	GasUsed uint64 `json:"gasUsed"          gencodec:"required"`
	//与此区块相关的附加数据。
	Extra []byte `json:"extraData"        gencodec:"required"`
	//区块哈希
	Hash common.Hash `json:"Hash"`
	//一个Hash值，当与nonce组合时，证明此区块已经执行了足够的计算。 nonce：POW生成的哈希值。
	MixDigest common.Hash `json:"mixHash"`
	Nonce     uint64      `json:"nonce"`

	Size common.StorageSize `json:"size"`
}
