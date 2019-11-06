package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/noah-blockchain/Auto-rewards/app"
	"github.com/noah-blockchain/Auto-rewards/config"
	"github.com/noah-blockchain/go-sdk/wallet"
	"github.com/robfig/cron/v3"
)

var (
	cfg = config.Config{}
)

func init() {
	minCoinDelegated := config.FloatGetEnv(os.Getenv("MIN_COINS_DELEGATED"), 1)
	cfg.StopListAccounts = strings.Split(os.Getenv("STOP_LIST"), ",")
	flag.StringVar(&cfg.SeedPhrase, "seed.phrase", os.Getenv("SEED_PHRASE"), "seed phrase not exist")
	flag.StringVar(&cfg.BaseCoin, "base.coin", os.Getenv("BASE_COIN"), "base coin not exist")
	flag.StringVar(&cfg.NodeApiURL, "node.api_url", os.Getenv("NODE_API_URL"), "node api url not exist")
	flag.StringVar(&cfg.ExplorerApiURL, "explorer.api_url", os.Getenv("EXPLORER_API_URL"), "explorer api url not exist")
	flag.StringVar(&cfg.Token, "token", os.Getenv("TOKEN"), "token not exist")
	flag.StringVar(&cfg.CronTime, "cron_time", os.Getenv("CRON_TIME"), "cron time not exist")
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
	case cfg.CronTime == "":
		log.Panicf("Invalid value %s for field %s", cfg.CronTime, "cron_time")
	}

	timeZoneMSK, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Panic(err)
	}

	cronScheduler := cron.New(
		cron.WithParser(
			cron.NewParser(
				cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow,
			),
		),
		cron.WithLogger(
			cron.VerbosePrintfLogger(
				log.New(os.Stdout,
					"go cron scheduler: ",
					log.LstdFlags|log.Lmicroseconds|log.Llongfile|log.Lshortfile,
				),
			),
		),
		cron.WithLocation(timeZoneMSK),
	)

	_, _ = cronScheduler.AddFunc(cfg.CronTime, func() {

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

		for _, item := range *multiSend {
			log.Println(fmt.Sprintf("Will be send %s %s to %s", item.Value.String(), item.Coin, item.To))
		}

		if err = autoRewards.SendMultiAccounts(walletFrom, *multiSend, "Auto-Reward payment", cfg.BaseCoin); err != nil {
			log.Panicln(err)
		}
		log.Println("All multi accounts was successful transferred.")

	})
	cronScheduler.Start()

	select {}
}
