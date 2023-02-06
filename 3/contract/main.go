package main

import (
	"context"
	"crypto/ecdsa"
	"eth-test/3/contract/store"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)           // in wei
	auth.GasLimit = uint64(300000 * 100) // in units 执行合约消耗的gas非常多, 所以这个需要设置的大一点, 合约执行失败的时候会显示需要多少gas
	auth.GasPrice = gasPrice

	input := "1.0"
	address, tx, instance, err := store.DeployStore(auth, client, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address.Hex())   // 0xAf6a1564d0C6F511409132F36E4f8124A4b0B8D2
	fmt.Println(tx.Hash().Hex()) // 0xbd0f996080294e28658061b110273a224ff0f603b16192f827f146bfdb91d936
	// GasLimit 设置的太小了, 所以下面的交易失败了
	// 0x3Ca429B885661Ed036b76799F50A9470AC59E684
	// 0x098f3d1f1d93b288f020a44796d78ddc5e098d0544c127f54f774c8213c0b943
	_ = instance
}
