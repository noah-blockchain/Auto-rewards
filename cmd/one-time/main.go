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
)

var cfg = config.Config{}

func init() {
	flag.StringVar(&cfg.SeedPhrase, "seed.phrase", os.Getenv("SEED_PHRASE"), "seed phrase not exist")
	flag.StringVar(&cfg.BaseCoin, "base.coin", os.Getenv("BASE_COIN"), "base coin not exist")
	flag.StringVar(&cfg.NodeApiURL, "node.api_url", os.Getenv("NODE_API_URL"), "node api url not exist")
	flag.StringVar(&cfg.ExplorerApiURL, "explorer.api_url", os.Getenv("EXPLORER_API_URL"), "explorer api url not exist")
	flag.StringVar(&cfg.Token, "token", os.Getenv("TOKEN"), "token not exist")
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

	autoRewards := app.NewAutoRewards(cfg)
	autoRewards.Task()
}
