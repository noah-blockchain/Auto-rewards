package config

import (
	"os"
	"strconv"
)

type Config struct {
	SeedPhrase       string
	BaseCoin         string
	NodeApiURL       string
	ExplorerApiURL   string
	Token            string
	CronTime         string
	MinCoinDelegated float64
	StopListAccounts []string
}

func FloatGetEnv(name string, def float64) float64 {
	value := os.Getenv(name)
	if value == "" {
		return def
	}
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return def
	}
	return v
}
