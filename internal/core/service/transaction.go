package service

import (
	"github.com/authorizer/internal/core/domain"
	"github.com/authorizer/internal/core/ports"
)

// Transaction service to process transactions
type Transaction struct {
	repo       ports.AccountRepository
	violations domain.Violations
}

// NewTransaction create a new Transaction instance
func NewTransaction(r ports.AccountRepository) Transaction {
	return Transaction{repo: r, violations: domain.Violations{}}
}

// Authorize process domain.Transaction and return domain.Account
func (t *Transaction) Authorize(transaction domain.Transaction) (*domain.Account, domain.Violations) {
	account, _ := t.repo.Retrieve(transaction.Time)

	if account == nil {
		return nil, domain.Violations{domain.AccountNotInitializedViolation}
	}

	if !account.Ledger.ActiveCard {
		return account, domain.Violations{domain.CardNotActiveViolation}
	}

	violations := t.validate(account, transaction)

	if len(violations) > 0 {
		return account, violations
	}

	account.Authorizations = append(account.Authorizations, domain.TransactionAuthorization{
		Merchant:       transaction.Merchant,
		Amount:         transaction.Amount,
		AvailableLimit: account.Ledger.AvailableLimit,
		Time:           transaction.Time,
	})

	changeAvailable(account, transaction.Amount)

	_ = t.repo.Update(*account)

	return account, violations
}

func (t *Transaction) validate(a *domain.Account, transaction domain.Transaction) domain.Violations {
	violations := domain.Violations{}

	violation := validateAvailable(a, transaction.Amount)
	violations.AddViolation(violation)

	rulesViolations := evaluateRules(a, transaction)
	violations.AddViolation(rulesViolations...)

	violation = validateIdempotencyTransactions(a.Authorizations, transaction)
	violations.AddViolation(violation)

	return violations
}

func validateAvailable(account *domain.Account, amount int64) string {
	if account.Ledger.AvailableLimit-amount < 0 {
		return domain.InsufficientLimitViolation
	}

	return ""
}

func changeAvailable(account *domain.Account, amount int64) {
	account.Ledger.AvailableLimit = account.Ledger.AvailableLimit - amount
}

func validateIdempotencyTransactions(authorizations []domain.TransactionAuthorization, t domain.Transaction) string {
	lenAuthorizations := len(authorizations)

	if lenAuthorizations == 0 {
		return ""
	}

	return validateDoubledTransactions(authorizations, t)
}

func validateDoubledTransactions(authorizations []domain.TransactionAuthorization, t domain.Transaction) string {
	lenAuthorizations := len(authorizations)

	for i := lenAuthorizations - 1; i >= 0; i-- {
		if t.Time.Sub(authorizations[i].Time).Minutes() > 2 {
			return ""
		}

		if t.Merchant == authorizations[i].Merchant && t.Amount == authorizations[i].Amount {
			return domain.DoubledTransactionViolation
		}
	}

	return ""
}

func evaluateRules(account *domain.Account, transaction domain.Transaction) []string {
	var violations []string

	for _, rule := range account.SpendingControl.Rules {
		violation := evaluateRule(&rule, transaction)

		if violation != "" {
			violations = append(violations, violation)
		}
	}

	return violations
}

func evaluateRule(rule *domain.Rule, transaction domain.Transaction) string {
	rule.Accumulator.AddSpend(transaction.Amount)

	if rule.Type == "usage-limit" && rule.Accumulator.CurrentPeriodUsed > rule.UsageLimit {
		return rule.RuleViolation
	}

	return ""
}
