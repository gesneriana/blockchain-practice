package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
)

func createKs() {
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	password := "gesneriana"
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex()) // 0xDe95E88936c39A5C81eb8fc101b3557231EFAFCf
}

func importKs() {
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)

	for _, account := range ks.Accounts() {
		fmt.Println(account.Address.Hex()) // 0xDe95E88936c39A5C81eb8fc101b3557231EFAFCf
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, errors.WithStack(err)
}

func main() {
	var tmpFolder, err = PathExists("./tmp")
	if !tmpFolder {
		if err == nil {
			createKs()
			return
		}
		log.Println(err)
		return
	}

	importKs()
}
