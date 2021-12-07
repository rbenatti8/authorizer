package domain

type SpendingControl struct {
	Rules []Rule
}

type Account struct {
	Ledger          Ledger
	SpendingControl SpendingControl
	Authorizations  []TransactionAuthorization
}
