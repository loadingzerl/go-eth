package modul

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
)

type Address struct {
	Address    common.Address `json:"address"`
	Balance    string         `json:"balance"`
	Code       abi.ABI        `json:"contractCode"`
	Nonce      uint64         `json:"nonce"`
	Number     *big.Int       `json:"blockHeight"`
	Time       uint64         `json:"timeStamp"`
	Storagekey string         `json:"stotageKey"`
}

func NewAddress() *Address {
	addressStatus := Address{
		Address:    common.HexToAddress(""),
		Balance:    "",
		Code:       abi.ABI{},
		Nonce:      0,
		Number:     big.NewInt(0),
		Time:       0,
		Storagekey: "",
	}
	return &addressStatus
}

func (a *Address) AddressStatus(block *types.Block, client *ethclient.Client) []*Address {
	addressStatus := make([]*Address, 0)
	for _, tx := range block.Transactions() {
		addressFrom := NewAddress()
		addressTo := NewAddress()
		msg, err := tx.AsMessage(types.LatestSignerForChainID(chainId), tx.GasPrice())
		if err == nil {
			addressFrom.Address = msg.From()
		}
		balanceFrom, _ := client.BalanceAt(context.Background(), addressFrom.Address, block.Number())
		addressFrom.Balance = balanceFrom.String()
		contractAbi, err := abi.JSON(strings.NewReader(TokenABI))
		if err != nil {
			fmt.Println(err)
		}
		addressFrom.Code = contractAbi
		addressFrom.Nonce = tx.Nonce()
		addressFrom.Number = block.Number()
		addressFrom.Time = block.Time()

		stotagekeyFrom, err := client.StorageAt(context.Background(), addressFrom.Address, tx.Hash(), block.Number())
		if err != nil {
			fmt.Println(err)
		}
		addressFrom.Storagekey = hex.EncodeToString(stotagekeyFrom)
		addressStatus = append(addressStatus, addressFrom)

		if tx.To() != nil {
			addressTo.Address = *tx.To()
		}
		balanceTo, _ := client.BalanceAt(context.Background(), addressFrom.Address, block.Number())
		addressTo.Balance = balanceTo.String()
		addressTo.Code = contractAbi
		addressTo.Nonce = tx.Nonce()
		addressTo.Number = block.Number()
		addressTo.Time = block.Time()
		stotagekeyTo, err := client.StorageAt(context.Background(), addressTo.Address, tx.Hash(), block.Number())
		if err != nil {
			fmt.Println(err)
		}
		addressTo.Storagekey = hex.EncodeToString(stotagekeyTo)
		addressStatus = append(addressStatus, addressTo)
	}
	return addressStatus
}
