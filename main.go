package main

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fbsobreira/gotron-sdk/pkg/keystore"
	todo "github.com/kjasuquo/todo-contracts/gen"
	"io/ioutil"
	"log"
	"math/big"
)

var infuraTestURL = "https://sepolia.infura.io/v3/8495aa8709f748338f32c1de7da0119e"

func main() {
	ctx := context.Background()
	b, err := ioutil.ReadFile("./wallet/UTC--2025-01-17T11-54-30.882022000Z--0ec715f34b8b32cc141fa83a4f2e6e1e494eed94")
	if err != nil {
		log.Fatal("cannot open keystore: ", err)
	}

	key, err := keystore.DecryptKey(b, "password")
	if err != nil {
		log.Fatal("cannot key from keystore: ", err)
	}

	client, err := ethclient.DialContext(ctx, infuraTestURL)
	if err != nil {
		log.Fatal("error connecting")
	}

	defer client.Close()

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal("cannot get suggested gas price ", err)
	}

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatal("chainID err: ", err)
	}

	newTodo, err := todo.NewTodo(common.HexToAddress("0x2439bD038aF926D31B2C5Bc17285cdc314005bb6"), client) // the contract address
	if err != nil {
		log.Fatal("chainID err: ", err)
	}

	opt, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, chainID)
	if err != nil {
		log.Fatal("opt err: ", err)
	}

	opt.GasPrice = gasPrice
	opt.GasLimit = uint64(3000000)
	//tx, err := newTodo.Add(opt, "first task")
	//if err != nil {
	//	log.Fatal("tx err: ", err)
	//}
	//
	//log.Println("tx hash.............................", tx.Hash()) // 0xd190e51670b415dcf4dac2ff11ebdd6383921215c034de7ec8ac7b3791ee5d53

	list, err := newTodo.List(&bind.CallOpts{
		From: crypto.PubkeyToAddress(key.PrivateKey.PublicKey),
	})
	if err != nil {
		log.Fatal("list err: ", err)
	}

	log.Println("list 1:.............................", list)

	update, err := newTodo.Update(opt, big.NewInt(0), "updated task")
	if err != nil {
		log.Fatal("update err: ", err)
	}

	log.Println("update hash.............................", update.Hash())

	list2, err := newTodo.List(&bind.CallOpts{
		From: crypto.PubkeyToAddress(key.PrivateKey.PublicKey),
	})
	if err != nil {
		log.Fatal("list err: ", err)
	}

	log.Println("list2:.............................", list2)

	toggle, err := newTodo.Toggle(opt, big.NewInt(0))
	if err != nil {
		log.Fatal("toggle err: ", err)
	}

	log.Println("toggle hash.............................", toggle.Hash())

	remove, err := newTodo.Remove(opt, big.NewInt(0))
	if err != nil {
		log.Fatal("remove err: ", err)
	}

	log.Println("toggle hash.............................", remove.Hash())
}
