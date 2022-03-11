package blockee

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Address struct {
	Address common.Address `json:"address"`
	Balance string `json:"balance"`
	Code abi.ABI `json:"contractCode"`
	Nonce uint64 `json:"nonce"`
	Number *big.Int `json:"blockHeight"`
	Time uint64 `json:"timeStamp"`
	Storagekey string `json:"stotageKey"`
}
