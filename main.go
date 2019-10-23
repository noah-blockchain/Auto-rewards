package main

import (
	"log"
	"os"

	"github.com/noah-blockchain/Auto-rewards/app"
	"github.com/noah-blockchain/go-sdk/wallet"
)

func main() {
	seed, _ := wallet.Seed(os.Getenv("SEED_PHRASE"))
	walletFrom, err := wallet.NewWallet(seed)
	if err != nil {
		log.Panicln(err)
	}

	multiSend, err := app.CreateMultiSendList(walletFrom.Address(), "NOAH")
	if err != nil || multiSend == nil {
		log.Panicln(err)
	}

	if err = app.SendMultiAccounts(walletFrom, *multiSend, "Payment from app", "NOAH"); err != nil {
		log.Panicln(err)
	}
}
