package repository

import (
	"github.com/authorizer/internal/core/domain"
	"github.com/authorizer/internal/driven/database"
	"github.com/authorizer/internal/dto"
	"time"
)

// AccountRepository represents a repository to domain.Account
type AccountRepository struct {
	db database.DB
}

// NewAccountRepository create a new AccountRepository instance
func NewAccountRepository(db database.DB) AccountRepository {
	return AccountRepository{db: db}
}

// Create insert new account on DB
func (ar AccountRepository) Create(account domain.Account) error {
	accDTO := buildDBEntity(account)
	_ = ar.db.Insert("accounts", 1, accDTO)
	return nil
}

// Update update account
func (ar AccountRepository) Update(account domain.Account) error {
	accDTO := buildDBEntity(account)
	_ = ar.db.Update("accounts", 1, accDTO)
	return nil
}

// Retrieve find account and return
func (ar AccountRepository) Retrieve(currentTime time.Time) (*domain.Account, error) {
	var accountDTO dto.Account
	err := ar.db.Find("accounts", 1, &accountDTO)

	if err != nil {
		return nil, err
	}

	return buildDomainAccount(currentTime, accountDTO), nil
}

func buildDBEntity(domainAccount domain.Account) dto.Account {
	transactionAuthorizations := make([]dto.TransactionAuthorization, 0, len(domainAccount.Authorizations))
	rules := make([]dto.Rule, 0, len(domainAccount.SpendingControl.Rules))

	for _, ta := range domainAccount.Authorizations {
		transactionAuthorizations = append(transactionAuthorizations, dto.TransactionAuthorization{
			Merchant:       ta.Merchant,
			Amount:         ta.Amount,
			AvailableLimit: ta.AvailableLimit,
			Time:           ta.Time,
		})
	}

	for _, rule := range domainAccount.SpendingControl.Rules {
		rules = append(rules, dto.Rule{
			Name:       rule.Name,
			Type:       rule.Type,
			UsageLimit: rule.UsageLimit,
			Accumulator: dto.Accumulator{
				Duration:           rule.Accumulator.Duration,
				CurrentPeriodUsed:  rule.Accumulator.CurrentPeriodUsed,
				CurrentPeriodSpend: rule.Accumulator.CurrentPeriodSpend,
				PeriodEndsDate:     rule.Accumulator.PeriodEndsDate,
			},
			RuleViolation: rule.RuleViolation,
		})
	}

	return dto.Account{
		Ledger: dto.Ledger{
			Active:         domainAccount.Ledger.ActiveCard,
			MaxLimit:       domainAccount.Ledger.MaxLimit,
			AvailableLimit: domainAccount.Ledger.AvailableLimit,
		},
		SpendingControl: dto.SpendingControl{Rules: rules},
		Transactions:    transactionAuthorizations,
	}
}

func buildDomainAccount(currentTime time.Time, accountDTO dto.Account) *domain.Account {
	transactionAuthorizations := make([]domain.TransactionAuthorization, 0, len(accountDTO.Transactions))
	rules := make([]domain.Rule, 0, len(accountDTO.SpendingControl.Rules))

	for _, ta := range accountDTO.Transactions {
		transactionAuthorizations = append(transactionAuthorizations, domain.TransactionAuthorization{
			Merchant:       ta.Merchant,
			Amount:         ta.Amount,
			AvailableLimit: ta.AvailableLimit,
			Time:           ta.Time,
		})
	}

	for _, rule := range accountDTO.SpendingControl.Rules {
		domainAccumulator := domain.BuildAccumulator(
			currentTime,
			rule.Accumulator.Duration,
			rule.Accumulator.CurrentPeriodUsed,
			rule.Accumulator.CurrentPeriodSpend,
			rule.Accumulator.PeriodEndsDate,
		)

		rules = append(rules, domain.Rule{
			Name:          rule.Name,
			Type:          rule.Type,
			UsageLimit:    rule.UsageLimit,
			Accumulator:   &domainAccumulator,
			RuleViolation: rule.RuleViolation,
		})
	}

	return &domain.Account{
		Ledger: domain.Ledger{
			ActiveCard:     accountDTO.Ledger.Active,
			MaxLimit:       accountDTO.Ledger.MaxLimit,
			AvailableLimit: accountDTO.Ledger.AvailableLimit,
		},
		SpendingControl: domain.SpendingControl{Rules: rules},
		Authorizations:  transactionAuthorizations,
	}
}
