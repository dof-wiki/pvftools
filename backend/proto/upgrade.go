package proto

type UpgradeItem struct {
	Level       int `json:"level"`
	FailRate    int `json:"fail_rate"`
	FailOp      int `json:"fail_op"`
	FailOpValue int `json:"fail_op_value"`
	CostId      int `json:"cost_id"`
	CostCount   int `json:"cost_count"`
}
