package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("wss://goerli.infura.io/ws/v3/a6fddc35a5ed4a4fa42e07392a969cd0")
	if err != nil {
		log.Fatal(err)
	}

	account := common.HexToAddress("0xDe95E88936c39A5C81eb8fc101b3557231EFAFCf")
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance) // 25893180161173005034

	blockNumber := big.NewInt(8375735)
	balanceAt, err := client.BalanceAt(context.Background(), account, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balanceAt) // 25729324269165216042

	fbalance := new(big.Float)
	fbalance.SetString(balanceAt.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Println(ethValue) // 25.729324269165216041

	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	fmt.Println(pendingBalance) // 25729324269165216042
	if err != nil {
		log.Fatal(err)
	}
}
