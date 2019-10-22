package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/noah-blockchain/Auto-rewards/models"
	"github.com/noah-blockchain/go-sdk/api"
	"github.com/noah-blockchain/go-sdk/transaction"
	"github.com/noah-blockchain/go-sdk/wallet"
)

func SendMultiAccounts(walletFrom *wallet.Wallet, dict []models.MultiSendItem, payload string, gasCoin string) error {
	nodeAPI := api.NewApi(os.Getenv("NODE_API_URL"))

	nonce, err := nodeAPI.Nonce(walletFrom.Address())
	if err != nil {
		return err
	}

	tx := transaction.NewMultisendData()
	for _, d := range dict {
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
		SetNonce(nonce).SetGasPrice(1).SetGasCoin(gasCoin).SetPayload([]byte(payload)).Sign(walletFrom.PrivateKey())
	if err != nil {
		return err
	}

	res, err := nodeAPI.Send(finishedTx)
	if err != nil {
		return err
	}

	if res.Error.Code != 0 {
		return errors.New(res.Error.Message)
	}

	fmt.Println(res.Result.Hash)
	fmt.Println(res.Result.Data)
	fmt.Println(res.Result.Log)

	return nil
}
