package main

import (
	"fmt"
	"os"

	"github.com/noah-blockchain/Auto-rewards/app"
	"github.com/noah-blockchain/go-sdk/wallet"
)

func main() {

	walletFrom, err := wallet.NewWallet([]byte(os.Getenv("SEED_PHRASE")))
	if err != nil {
		fmt.Println(err)
		return
	}

	multiSend, err := app.CreateMultiSendList(walletFrom.Address(), "NOAH")
	if err != nil || multiSend == nil {
		fmt.Println(err)
		return
	}

	if err = app.SendMultiAccounts(walletFrom, *multiSend, "Payment from app", "NOAH"); err != nil {
		fmt.Println(err)
		return
	}
}
