package dto

import "time"

type TransactionOperation struct {
	Merchant string    `json:"merchant"`
	Amount   int64     `json:"amount"`
	Time     time.Time `json:"time"`
}
