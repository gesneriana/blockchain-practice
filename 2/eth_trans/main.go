package main

import (
	"context"
	"crypto/ecdsa"
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

	// 可以使用 MetaMask 钱包导出私钥
	privateKey, err := crypto.HexToECDSA("f2e1027fe779a0b83966cae2d0a083cefd8c53ff63562a407aa1a686bb13f918")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	// 在生产环境中千万不要这样获取钱包地址,  钱包地址直接从数据库获取更安全, 从私钥生成钱包地址需要读取私钥, 能不读取私钥就不要直接读取私钥
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// https://twitter.com/asuna_kizuna/status/1619431347466166273?s=20&t=9iWW4qcgck1olk3EnqPjjw

	value := big.NewInt(1000000000000000000 / 10000) // in wei (1/10000 eth)
	gasLimit := uint64(21000)                        // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 接收人的钱包地址
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

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	// 交易hash 0x10447a1d9b39bbddb38b54729061c1a806c337b9f35b5e37c21cc97cc7808727
	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}
