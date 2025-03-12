package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fbsobreira/gotron-sdk/pkg/keystore"
	todo "github.com/kjasuquo/todo-contracts/gen"
)

var infuraTestURL = "https://sepolia.infura.io/v3/8495aa8709f748338f32c1de7da0119e"

func main() {
	log.Println("start...................................................")

	b, err := ioutil.ReadFile("./wallet/UTC--2025-01-17T11-54-30.882022000Z--0ec715f34b8b32cc141fa83a4f2e6e1e494eed94")
	if err != nil {
		log.Fatal("cannot open keystore: ", err)
	}

	key, err := keystore.DecryptKey(b, "password")
	if err != nil {
		log.Fatal("cannot key from keystore: ", err)
	}
	ctx := context.Background()

	client, err := ethclient.DialContext(ctx, infuraTestURL)
	if err != nil {
		log.Fatal("error connecting")
	}

	add := crypto.PubkeyToAddress(key.PrivateKey.PublicKey)

	bal, err := client.BalanceAt(ctx, add, nil)
	if err != nil {
		log.Fatal("error getting balance2: ", err)
	}

	log.Println("the address balance before deployment.........................: ", bal)

	pendingNonce, err := client.PendingNonceAt(ctx, add)
	if err != nil {
		log.Fatal("cannot get pending nonce ", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal("cannot get suggested gas price ", err)
	}

	log.Println("the gasPrice.........................: ", gasPrice)

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatal("chainID err: ", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, chainID)
	if err != nil {
		log.Fatal("auth err: ", err)
	}
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(3000000)
	auth.Nonce = big.NewInt(int64(pendingNonce))

	// todo.DeployKJW3Token(auth, client, "KJToken" , "KJT") for deploying token

	a, tx, _, err := todo.DeployTodo(auth, client)
	if err != nil {
		log.Fatal("deployTodo err: ", err)
	}

	log.Println("the address...............", a.Hex())
	log.Println("the tx hash...............", tx.Hash().Hex())

	// 0x2439bD038aF926D31B2C5Bc17285cdc314005bb6 - contract address
	// 0xc777137cdc509d7e86d48b71c8a91a00d894aabbe03de243756b3ad4b0a36792 - contract hash

	log.Println("end...................................................")
}
