package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("wss://goerli.infura.io/ws/v3/a6fddc35a5ed4a4fa42e07392a969cd0")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("f2e1027fe779a0b83966cae2d0a083cefd8c53ff63562a407aa1a686bb13f918")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1000000000000000000 / 10000) // in wei (1/10000 eth)
	gasLimit := uint64(21000)                        // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0xf0046E53d11c6E7a6BAdB4C9E849fCC26509dc16")
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	defer buf.Reset()
	err = signedTx.EncodeRLP(buf)
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个未发布的交易, 这个场景不常用,  比如钱包里面没有gas了, 交易失败了,  所以发不出去, 可以先创建交易,  等有了gas再发布
	fmt.Printf("buf hex: %s\n", hex.EncodeToString(buf.Bytes()))
}
