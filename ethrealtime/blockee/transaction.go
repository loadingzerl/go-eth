package blockee

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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
	Value             string             `json:"txnValue"`
	Logs              []*Log             `json:"logs"`
	ContractAddress   common.Address     `json:"txnContractAddress"`
	Size              common.StorageSize `json:"size"`
	TxnType           uint8              `json:"txnType"`
	TxnStatus         uint64             `json:"txnStatus"`
	ChainId           *big.Int           `json:"chainId"`
	PostState         string             `json:"txnPostState"`
	MaxPriority       *big.Int           `json:"txnMaxPriority"` //gasprice - GasFeeCap
	GasTipCap         *big.Int           `json:"txnGasTipCap"`
	GasFeeCap         *big.Int           `json:"txnGasFeeCap"`
	CumulativeGasUsed uint64             `json:"txnCumulativeGasUsed"`
	TransactionIndex  uint               `json:"txnTransactionIndex"`
	AccessList        types.AccessList   `json:"txnAccessList"`
	TxnGasUsed        uint64             `json:"txGasUsed"`

	GasMaxFee *big.Int `json:"txnMaxFee"`
}

func NewTransaction() *Transaction {
	transaction := Transaction{
		Hash:              common.HexToHash(""),
		From:              common.HexToAddress(""),
		To:                nil,
		Gas:               0,
		GasPrice:          big.NewInt(0),
		Data:              "",
		Nonce:             0,
		Value:             "",
		Logs:              nil,
		ContractAddress:   common.HexToAddress(""),
		Size:              0,
		TxnType:           0,
		TxnStatus:         0,
		ChainId:           big.NewInt(1),
		PostState:         "",
		MaxPriority:       big.NewInt(0),
		GasTipCap:         big.NewInt(0),
		GasFeeCap:         big.NewInt(0),
		CumulativeGasUsed: 0,
		TransactionIndex:  0,
		AccessList:        nil,
		TxnGasUsed:        0,
		GasMaxFee:         big.NewInt(0),
	}
	return &transaction
}

func (t *Transaction) Transactions(block *types.Block, client *ethclient.Client) []*Transaction {
	transactions := make([]*Transaction, 0)
	var chainId = new(big.Int).SetUint64(uint64(1))
	for _, tx := range block.Transactions() {
		t.Hash = tx.Hash()
		msg, err := tx.AsMessage(types.LatestSignerForChainID(chainId), tx.GasPrice())
		if err == nil {
			t.From = msg.From()
		}
		t.To = tx.To()
		if t.To == nil {
			to := crypto.CreateAddress(msg.From(), msg.Nonce())
			t.ContractAddress = to
		}
		t.Gas = tx.Gas()
		t.Data = hex.EncodeToString(tx.Data())
		t.Nonce = tx.Nonce()
		t.Value = tx.Value().String()
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			fmt.Println("tx receipt ", receipt)
		}
		if receipt != nil {
			t.ContractAddress = receipt.ContractAddress
			t.TxnType = receipt.Type
			t.TxnStatus = receipt.Status
			t.PostState = hex.EncodeToString(receipt.PostState)
			t.TxnGasUsed = receipt.GasUsed
			t.CumulativeGasUsed = receipt.CumulativeGasUsed
			t.TransactionIndex = receipt.TransactionIndex
			t.PostState = hex.EncodeToString(receipt.PostState)
		}
		t.Size = tx.Size()
		t.ChainId = tx.ChainId()

		t.GasMaxFee = tx.GasFeeCap()
		t.GasFeeCap = block.Header().BaseFee

		t.GasPrice = tx.GasPrice()
		t.MaxPriority = tx.GasTipCap()
		t.GasTipCap = tx.GasTipCap()
		if t.TxnType == 2 {
			price := big.Int{}
			if price.Add(t.GasFeeCap, t.MaxPriority).Cmp(t.GasMaxFee) <= 0 {
				t.GasPrice = price.Add(t.GasFeeCap, t.MaxPriority)
			} else {
				t.GasPrice = t.GasMaxFee
			}
		}
		t.AccessList = tx.AccessList()

		transactions = append(transactions, t)
	}
	return transactions
}
