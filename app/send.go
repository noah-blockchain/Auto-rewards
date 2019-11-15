package app

import (
	"errors"
	"fmt"
	"log"

	"github.com/noah-blockchain/Auto-rewards/models"
	"github.com/noah-blockchain/go-sdk/api"
	"github.com/noah-blockchain/go-sdk/transaction"
	"github.com/noah-blockchain/go-sdk/wallet"
)

func (a AutoRewards) SendMultiAccounts(walletFrom *wallet.Wallet, txs []models.MultiSendItem, payload string, gasCoin string) error {
	if len(txs) == 0 {
		fmt.Println("ERROR! Empty txs list")
		return errors.New("ERROR! Multi list accounts cant be empty")
	}

	nodeAPI := api.NewApi(a.cfg.NodeApiURL)

	nonce, err := nodeAPI.Nonce(walletFrom.Address())
	if err != nil {
		return err
	}

	tx := transaction.NewMultisendData()
	for _, d := range txs {
		tx.AddItem(
			*transaction.NewMultisendDataItem().
				SetCoin(d.Coin).
				SetValue(d.Value).
				MustSetTo(d.To),
		)
	}

	signedTx, err := transaction.NewBuilder(transaction.TestNetChainID).NewTransaction(tx)
	if err != nil {
		return err
	}

	finishedTx, err := signedTx.
		SetNonce(nonce).SetGasPrice(255).SetGasCoin(gasCoin).SetPayload([]byte(payload)).Sign(walletFrom.PrivateKey())
	if err != nil {
		return err
	}

	res, err := nodeAPI.SendTransaction(finishedTx)
	if err != nil {
		return err
	}

	log.Println("OK! MultiSend trx successful created with HASH=", res.Hash)
	return nil
}
