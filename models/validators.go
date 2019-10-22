package models

type ValidatorItem struct {
	PubKey      string `json:"pub_key"`
	VotingPower string `json:"voting_power"`
}

type ValidatorList struct {
	Validators []ValidatorItem `json:"result"`
}
