package models

type BalanceItem struct {
	Coin   string `json:"coin"`
	Amount string `json:"amount"`
}
type BalancesList struct {
	Balances []BalanceItem `json:"balances"`
}

type AddressInfo struct {
	Data BalancesList `json:"data"`
}
