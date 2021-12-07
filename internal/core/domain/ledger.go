package domain

type Ledger struct {
	ActiveCard     bool
	MaxLimit       int64
	AvailableLimit int64
}
