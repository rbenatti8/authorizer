package service

import (
	"time"

	"github.com/authorizer/internal/core/domain"
	"github.com/authorizer/internal/core/ports"
)

type Account struct {
	repo ports.AccountRepository
}

func NewAccount(r ports.AccountRepository) Account {
	return Account{repo: r}
}
func (a Account) InitAccount(activeCard bool, maxLimit int64) (*domain.Account, []string) {
	existentAccount, _ := a.repo.Retrieve(time.Now())

	if existentAccount != nil {
		return existentAccount, []string{domain.AccountAlreadyInitializedViolation}
	}

	newAccount := domain.Account{
		Ledger: domain.Ledger{
			ActiveCard:     activeCard,
			MaxLimit:       maxLimit,
			AvailableLimit: maxLimit,
		},
		SpendingControl: domain.SpendingControl{
			Rules: []domain.Rule{
				{
					Name:       "max transactions in 2 minutes",
					Type:       "usage-limit",
					UsageLimit: 3,
					Accumulator: &domain.Accumulator{
						Duration:          2 * time.Minute,
						CurrentPeriodUsed: 0,
					},
					RuleViolation: "high-frequency-small-interval",
				},
			},
		},
		Authorizations: []domain.TransactionAuthorization{},
	}

	_ = a.repo.Create(newAccount)

	return &newAccount, []string{}
}
