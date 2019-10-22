package models

type DelegatorItem struct {
	Coin      string `json:"coin"`
	Address   string `json:"address"`
	Value     string `json:"value"`
	NoahValue string `json:"noah_value"`
}

type ValidatorData struct {
	DelegatorCount int64           `json:"delegator_count"`
	DelegatorList  []DelegatorItem `json:"delegator_list"`
}

type ValidatorInfo struct {
	Result ValidatorData `json:"data"`
}
