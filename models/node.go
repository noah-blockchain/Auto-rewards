package models

type NodeResult struct {
	LatestBlockHeight string `json:"latest_block_height"`
}

type NodeStatus struct {
	Result NodeResult `json:"result"`
}
