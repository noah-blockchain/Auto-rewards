package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/noah-blockchain/Auto-rewards/app"
	"github.com/noah-blockchain/Auto-rewards/config"
	"github.com/noah-blockchain/go-sdk/wallet"
)

var cfg = config.Config{}

func init() {
	minCoinDelegated := config.FloatGetEnv(os.Getenv("MIN_COINS_DELEGATED"), 100)
	cfg.StopListAccounts = strings.Split(os.Getenv("STOP_LIST"), ",")
	flag.StringVar(&cfg.SeedPhrase, "seed.phrase", os.Getenv("SEED_PHRASE"), "seed phrase not exist")
	flag.StringVar(&cfg.BaseCoin, "base.coin", os.Getenv("BASE_COIN"), "base coin not exist")
	flag.StringVar(&cfg.NodeApiURL, "node.api_url", os.Getenv("NODE_API_URL"), "node api url not exist")
	flag.StringVar(&cfg.ExplorerApiURL, "explorer.api_url", os.Getenv("EXPLORER_API_URL"), "explorer api url not exist")
	flag.StringVar(&cfg.Token, "token", os.Getenv("TOKEN"), "token not exist")
	flag.Float64Var(&cfg.MinCoinDelegated, "min_coin_delegated", minCoinDelegated, "min coin delegated not setup")
}

func main() {
	flag.Usage = func() {
		cmd := strings.TrimSuffix(path.Base(os.Args[0]), ".test")
		fmt.Printf("Usage of %s:\n", cmd)
		flag.PrintDefaults()
	}
	flag.Parse()

	switch {
	case cfg.SeedPhrase == "":
		log.Panicf("Invalid value %s for field %s", cfg.SeedPhrase, "seed.phrase")
	case cfg.BaseCoin == "":
		log.Panicf("Invalid value %s for field %s", cfg.BaseCoin, "base.coin")
	case cfg.NodeApiURL == "":
		log.Panicf("Invalid value %s for field %s", cfg.NodeApiURL, "node.api_url")
	case cfg.ExplorerApiURL == "":
		log.Panicf("Invalid value %s for field %s", cfg.ExplorerApiURL, "explorer.api_url")
	case cfg.Token == "":
		log.Panicf("Invalid value %s for field %s", cfg.Token, "token")
	}

	seed, _ := wallet.Seed(cfg.SeedPhrase)
	walletFrom, err := wallet.NewWallet(seed)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Wallet was successful received.")

	autoRewards := app.NewAutoRewards(cfg)
	multiSend, err := autoRewards.CreateMultiSendList(walletFrom.Address(), cfg.BaseCoin)
	if err != nil || multiSend == nil {
		log.Panicln(err)
	}
	log.Println("Multi list for accounts was successful created.")

	if err = autoRewards.SendMultiAccounts(walletFrom, *multiSend, "Payment from app", cfg.BaseCoin); err != nil {
		log.Panicln(err)
	}
	log.Println("All multi accounts was successful transferred.")
}
