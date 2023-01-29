package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"eth-test/2/erc20_trans/token"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
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

	// 发送人的钱包地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// 接收人的钱包地址
	toAddress := common.HexToAddress("0x4281ecf07378ee595c564a59048801330f3084ee")
	// ChainLink 代币合约地址
	tokenAddress := common.HexToAddress("0x326c977e6efc84e512bb9c30f76e30c160ed06fb")
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	amount := new(big.Int)
	amount.SetString("1000000000000000000", 10) // 1 tokens
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 函数方法切片传递
	transferFnSignature := []byte("transfer(address,uint256)")
	// 生成签名hash
	hash := sha3.New384()
	hash.Write(transferFnSignature)
	// 切片
	methID := hash.Sum(nil)[:4]
	// 填充地址
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	// 代币量填充
	var data []byte
	data = append(data, methID...)
	data = append(data, paddedAddress...)
	data = append(data, common.LeftPadBytes(amount.Bytes(), 32)...)

	// 使用方法估算燃气费
	gasPrice, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To:   &toAddress,
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}

	var opts = &bind.TransactOpts{
		From:  fromAddress,
		Nonce: big.NewInt(int64(nonce)),
		Signer: func(senderAddress common.Address, trans *types.Transaction) (*types.Transaction, error) {
			keyAddr := crypto.PubkeyToAddress(*publicKeyECDSA)
			if senderAddress != keyAddr {
				return nil, errors.New("not authorized to sign this account")
			}
			return types.SignTx(trans, types.NewEIP155Signer(chainID), privateKey)
		},
		Value:    big.NewInt(0),               // 这里是代币转账, 所以ETH填0
		GasPrice: big.NewInt(int64(gasPrice)), // gas 估算, 可能不准确, 需要在主链上测试才知道, 或者先查看其他的智能合约大概需要多少gas作为参考
		GasLimit: 3000000,                     // 推荐设置, 如果一次转账给多个人可能会消耗大量gas
		Context:  context.Background(),
	}
	signedTx, err := instance.Transfer(opts, toAddress, amount)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex()) // tx sent: 0xc399a3c88924301c2b4091ed58961d5b8e6cf7f3dc2f8c90974dc687a939eaf7
}
