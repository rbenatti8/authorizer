package domain

import "time"

type TransactionAuthorization struct {
	Merchant       string
	Amount         int64
	AvailableLimit int64
	Time           time.Time
}
