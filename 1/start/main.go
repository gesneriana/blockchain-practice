package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// cloudflare-eth 是以太坊主网的网关, 这个是生产环境的网络
	// https://cloudflare-eth.com

	// mainnet.infura.io 也是eth的主网
	// wss://mainnet.infura.io/ws/v3/a6fddc35a5ed4a4fa42e07392a969cd0

	// goerli 是测试网络的url
	// wss://goerli.infura.io/ws/v3/a6fddc35a5ed4a4fa42e07392a969cd0
	client, err := ethclient.Dial("wss://goerli.infura.io/ws/v3/a6fddc35a5ed4a4fa42e07392a969cd0")
	if err != nil {
		log.Fatal(err)
	}

	chainId, _ := client.ChainID(context.Background())
	netId, _ := client.NetworkID(context.Background())
	fmt.Printf("chainId: %v, netId: %v\n", chainId, netId)
}
