package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

func main() {
	client, err := ethclient.Dial("wss://goerli.infura.io/ws/v3/a6fddc35a5ed4a4fa42e07392a969cd0")
	if err != nil {
		log.Fatal(err)
	}

	rawTx := "f865021682520894f0046e53d11c6e7a6badb4c9e849fcc26509dc16865af3107a4000802ea00e77e734f3af5f285e949a03504fb1e9b417f0b4567944412ca0dda9a15302cba07c989cf44c74facbfe60db92309f4c2a17566fde88f7b992b757c33dd4271e1b"

	rawTxBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		log.Fatal(err)
	}

	tx := new(types.Transaction)
	err = rlp.DecodeBytes(rawTxBytes, tx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", tx.Hash().Hex()) // tx sent: 0x53e1a2388084969bf72b2759f303529558fb55a407cf888f47ada92f6e6d6f67
}
