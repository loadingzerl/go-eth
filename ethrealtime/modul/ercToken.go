package modul

import (
	"context"
	"ethernum/token"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type ERCToken struct {
	//转账人
	From common.Address `json:"from"`
	//接收人
	To common.Address `json:"to"`
	//交易代币金额
	Value string `json:"value"`
	//创建合约的地址 from
	ContractAddress string `json:"creator"`
	//代币交易哈希
	Hash common.Hash `json:"hash"`
	//当前区块from的余额
	FromBalance string `json:"fromBalance"`
	//当前区块to的余额
	ToBalance string `json:"toBalance"`
	//代币名全称
	Name string `json:"name"`
	//代币缩写
	Symbol string `json:"symbol"`

	//合约地址
	Address        common.Address `json:"address"`
	BlockHash      common.Hash    `json:"blockHash"`
	BlockNumber    *big.Int       `json:"blockNumber"`
	Decimals       uint8          `json:"decimals"`
	ERCType        string         `json:"ercType"`
	Time           uint64         `json:"timeStamp"`
	AddressBalance string         `json:"addressBalance"`
}

func NewERCToken() *ERCToken {
	erctoken := ERCToken{
		From:            common.HexToAddress(""),
		To:              common.HexToAddress(""),
		Value:           "",
		ContractAddress: "",
		Hash:            common.HexToHash(""),
		FromBalance:     "",
		ToBalance:       "",
		Name:            "",
		Symbol:          "",
		Address:         common.HexToAddress(""),
		BlockHash:       common.HexToHash(""),
		BlockNumber:     big.NewInt(0),
		Decimals:        0,
		ERCType:         "",
		Time:            0,
		AddressBalance:  "",
	}
	return &erctoken
}

func (ercToken *ERCToken) ERCToken(log *types.Log, client *ethclient.Client, topices []common.Hash, contractAbi abi.ABI, tx *types.Transaction, instance *token.Token) {
	switch topices[0] {

	case Transferevent:
		event, err := contractAbi.Unpack("Transfer", log.Data)
		if err != nil {
			//log.Println("event err:", err)
			fmt.Println("event err:", err)
		}
		//ercToken := modul.ERCToken{}
		if len(event) > 0 {
			//	"value:", event[0], "contractAddr:", v.Address, "TXhash:", tx.Hash())
			var vl *big.Int
			if value, ok := event[0].(*big.Int); ok {
				vl = value
			}
			if len(log.Topics) >= 3 {
				ercToken.From = common.HexToAddress(log.Topics[1].Hex())
				ercToken.To = common.HexToAddress(log.Topics[2].Hex())
			}

			//ercToken.ContractAddress = receipt.ContractAddress.String()
			//if receipt.ContractAddress == common.HexToAddress("0x0000000000000000000000000000000000000000"){
			//	ercToken.ContractAddress = ""
			//}

			if tx.To() == nil {
				msg, err := tx.AsMessage(types.LatestSignerForChainID(chainId), tx.GasPrice())
				if err != nil {
					fmt.Printf("error %s\n", err.Error())
				}
				//合约创建地址
				ercToken.ContractAddress = msg.From().String()
			}

			ercToken.Value = vl.String()
			ercToken.Address = log.Address
			ercToken.Hash = tx.Hash()
		} else {
			if len(log.Topics) >= 3 {
				ercToken.From = common.HexToAddress(log.Topics[1].Hex())
				ercToken.To = common.HexToAddress(log.Topics[2].Hex())
			}

			ercToken.Address = log.Address

			//ercToken.ContractAddress = receipt.ContractAddress.String()
			//if receipt.ContractAddress == common.HexToAddress("0x0000000000000000000000000000000000000000"){
			//	ercToken.ContractAddress = ""
			//}

			if tx.To() == nil {
				msg, err := tx.AsMessage(types.LatestSignerForChainID(chainId), tx.GasPrice())
				if err != nil {
					fmt.Printf("error %s\n", err.Error())
				}
				//合约创建地址
				ercToken.ContractAddress = msg.From().String()
			}

			ercToken.Hash = tx.Hash()
		}
		name, err := instance.Name(nil)
		if err != nil {
			fmt.Println("ERCTOKENname err:", err)
		}
		Symbol, err := instance.Symbol(nil)
		if err != nil {
			fmt.Println("ERCTOKENSymbol err:", err)
		}
		addressBalance, err := client.BalanceAt(context.Background(), ercToken.Address, nil)
		if err != nil {
			fmt.Println(err)
		}
		ercToken.AddressBalance = addressBalance.String()
		var fromBalance *big.Int
		var toBalance *big.Int
		if len(log.Topics) >= 3 {
			fromBalance, err = instance.BalanceOf(nil, common.HexToAddress(log.Topics[1].Hex()))
			if err != nil {
				fmt.Println("fromBalance", err)
			}
			toBalance, err = instance.BalanceOf(nil, common.HexToAddress(log.Topics[2].Hex()))
			if err != nil {
				fmt.Println("toBalance", err)
			}
		}
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())

		ercToken.BlockNumber = receipt.BlockNumber
		ercToken.BlockHash = receipt.BlockHash
		decimals, err := instance.Decimals(nil)
		if err != nil {
			fmt.Println("decimals err ", err)
		}

		ercToken.Time = 1

		ercToken.Decimals = decimals
		ercToken.FromBalance = fromBalance.String()
		ercToken.ToBalance = toBalance.String()
		ercToken.Name = name
		ercToken.Symbol = Symbol
		//state := TokenType(ercToken.Address)
		ercToken.ERCType = "ERC20"
	}
}

func (ercToken *ERCToken) TokenType(contractAddress common.Address, client *ethclient.Client) bool {
	erc721contract, err := token.NewErc721(contractAddress, client)
	if err != nil {
		fmt.Println(err)
	}
	data := [4]byte{0x80, 0xac, 0x58, 0xcd}
	iserc721, err := erc721contract.SupportsInterface(nil, data)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("查询ERC721结果:",iserc721)
	return iserc721
}
