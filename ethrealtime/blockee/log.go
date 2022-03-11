package blockee

import "github.com/ethereum/go-ethereum/common"

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
	Index uint `json:"logIndex"`

	Time uint64 `json:"timeStamp"`

	// The Removed field is true if this log was reverted due to a chain reorganisation.
	// You must pay attention to this field if you receive logs through a filter query.
	Removed bool `json:"removed"`

	ERCTokens []ERCToken `json:"erctokens"`
	ERCType string `json:"ercType"`


}
