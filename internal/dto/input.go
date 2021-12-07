package dto

type Input struct {
	Account     *AccountOperation     `json:"account,omitempty"`
	Transaction *TransactionOperation `json:"transaction,omitempty"`
}
