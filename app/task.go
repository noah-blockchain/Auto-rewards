package app

import (
	"fmt"
	"log"
	"time"

	"github.com/noah-blockchain/Auto-rewards/models"
	"github.com/noah-blockchain/go-sdk/api"
	"github.com/noah-blockchain/go-sdk/wallet"
)

const (
	WrongNonce        = 101
	InsufficientFunds = 107
)

func (a AutoRewards) Task() {
	seed, err := wallet.Seed(a.cfg.SeedPhrase)
	if err != nil {
		log.Panicln(err)
	}

	walletFrom, err := wallet.NewWallet(seed)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("OK! Wallet was successful received.")

	var multiSend []models.MultiSendItem
	attempt := 1
	for {
		fmt.Println("INFO! CreateMultiSendList Attempt number", attempt)

		if len(multiSend) == 0 {
			multiSend, err = a.CreateMultiSendList(walletFrom.Address(), a.cfg.BaseCoin)
			if err != nil || len(multiSend) == 0 {
				log.Println("ERROR! Multi send list not created", err)
				time.Sleep(15 * time.Second)
				attempt++
				continue
			}
			log.Println("OK! Multi list for accounts was successful created.")
		}

		err = a.SendMultiAccounts(walletFrom, multiSend, "Auto-Reward payment", a.cfg.BaseCoin)
		if err != nil {
			eTxError, ok := err.(*api.TxError)
			if ok && (eTxError.TxResult.Code != InsufficientFunds && eTxError.TxResult.Code != WrongNonce) {
				log.Println("SYSTEM ERROR!", err)
				break
			}

			log.Println("ERROR! Multi send list cant be send", err)
			time.Sleep(15 * time.Second)
			attempt++
			continue
		}
		log.Println("OK! Multi send list successful send")
		break
	}
}
