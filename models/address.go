package models

type BalanceItem struct {
	Coin   string `json:"coin"`
	Amount string `json:"amount"`
}
type BalanceInfo struct {
	Address  string        `json:"address"`
	Balances []BalanceItem `json:"balances"`
}

type AddressInfo struct {
	Data BalanceInfo `json:"data"`
}
