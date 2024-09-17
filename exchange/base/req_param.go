package base

type TransferParam struct {
	FromUser     string  `json:"from_user"`
	FromAccount  string  `json:"from_account"`
	FromType     string  `json:"from_type"`
	ToUser       string  `json:"to_user"`
	ToAccount    string  `json:"to_account"`
	ToType       string  `json:"to_type"`
	Assert       string  `json:"assert"`
	Amount       float64 `json:"amount"`
	TransferType string  `json:"transfer_type"`
}
