package dto

import "time"

type Accumulator struct {
	Duration           time.Duration `json:"duration"`
	CurrentPeriodUsed  int64         `json:"current_period_used"`
	CurrentPeriodSpend int64         `json:"current_period_spend"`
	PeriodEndsDate     time.Time     `json:"period_ends_date"`
}

type Rule struct {
	Name          string      `json:"name"`
	Type          string      `json:"type"`
	UsageLimit    int64       `json:"usage_limit"`
	Accumulator   Accumulator `json:"accumulator"`
	RuleViolation string      `json:"rule_violation"`
}

type SpendingControl struct {
	Rules []Rule `json:"rules"`
}

type Ledger struct {
	Active         bool  `json:"active"`
	MaxLimit       int64 `json:"max_limit"`
	AvailableLimit int64 `json:"available_limit"`
}

type TransactionAuthorization struct {
	Merchant       string    `json:"merchant"`
	Amount         int64     `json:"amount"`
	AvailableLimit int64     `json:"available_limit"`
	Time           time.Time `json:"time"`
}

type Account struct {
	Ledger          Ledger                     `json:"ledger"`
	SpendingControl SpendingControl            `json:"spending_control"`
	Transactions    []TransactionAuthorization `json:"transactions"`
}
