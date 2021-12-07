package domain

import "time"

type Transaction struct {
	Merchant string
	Amount   int64
	Time     time.Time
}
