package modul

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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
	Transactions []*Transaction
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
	Nonce uint64 `json:"blockNonce"`
	//Size
	Size common.StorageSize `json:"blockSize"`

	ReceivedFrom interface{}     `json:"receivedFrom"`
	UnclesNumber int             `json:"unclesNumber"`
	ReceivedAa   time.Time       `json:"receivedAt"`
	Coinbase     common.Address  `json:"miner"            gencodec:"required"`
	Uncles       []*types.Header `json:"uncles"`
	BlockBaseFee *big.Int        `json:"blockBaseFee"`
}

func NewBlock() *Block {
	block := Block{
		Number:       big.NewInt(0),
		Time:         0,
		ParentHash:   common.HexToHash(""),
		UncleHash:    common.HexToHash(""),
		Root:         common.HexToHash(""),
		TxHash:       common.HexToHash(""),
		ReceiptHash:  common.HexToHash(""),
		Difficulty:   big.NewInt(0),
		GasLimit:     0,
		GasUsed:      0,
		Extra:        "",
		Hash:         common.HexToHash(""),
		MixDigest:    common.HexToHash(""),
		Nonce:        0,
		Size:         0,
		ReceivedFrom: nil,
		UnclesNumber: 0,
		ReceivedAa:   time.Now(),
		Coinbase:     common.HexToAddress(""),
		Uncles:       nil,
		BlockBaseFee: big.NewInt(0),
		Transactions: nil,
	}
	return &block
}

func (b *Block) BlockObtain(block *types.Block, client *ethclient.Client) {
	//blockInit:
	//	block, err := client.BlockByNumber(context.Background(), number)
	//	if err != nil {
	//		fmt.Println("BlockNumber() err : ", err)
	//		clientInit()
	//		//block, err = client.BlockByNumber(context.Background(), number)
	//		//goto blockInit
	//	}
	//创建block对象
	b.Number = block.Number()
	b.Time = block.Time()
	b.ParentHash = block.ParentHash()
	b.UncleHash = block.UncleHash()
	b.Root = block.Root()
	b.TxHash = block.TxHash()
	b.ReceiptHash = block.ReceiptHash()
	b.Difficulty = block.Difficulty()
	b.GasLimit = block.GasLimit()
	b.GasUsed = block.GasUsed()
	b.Extra = hex.EncodeToString(block.Extra())
	b.Hash = block.Hash()
	b.MixDigest = block.MixDigest()
	b.Nonce = block.Nonce()
	b.Size = block.Size()
	b.ReceivedFrom = block.ReceivedFrom
	b.UnclesNumber = len(block.Uncles())
	b.ReceivedAa = block.ReceivedAt
	b.Coinbase = block.Coinbase()
	b.Uncles = block.Uncles()
	b.BlockBaseFee = block.BaseFee()

	//遍历交易信息t
	//transactions := Transactions(block)
	transaction := NewTransaction()
	b.Transactions = transaction.Transactions(block, client)

}
