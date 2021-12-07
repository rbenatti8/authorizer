package service

import (
	"github.com/authorizer/internal/core/domain"
	"github.com/authorizer/internal/driven/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAccount_InitAccount_Without_Violations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedAccount := domain.Account{
		Ledger: domain.Ledger{
			ActiveCard:     true,
			MaxLimit:       200,
			AvailableLimit: 200,
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

	accountRepoMock := repository.NewMockAccountRepository(ctrl)

	accountRepoMock.EXPECT().Retrieve(gomock.Any()).Return(nil, nil)
	accountRepoMock.EXPECT().Create(expectedAccount).Return(nil)

	as := NewAccount(accountRepoMock)

	account, _ := as.InitAccount(true, 200)

	assert.Equal(t, account.Ledger.AvailableLimit, int64(200))
	assert.Equal(t, account.Ledger.ActiveCard, true)
}

func TestAccount_InitAccount_With_Violations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccount := domain.Account{
		Ledger: domain.Ledger{
			ActiveCard:     true,
			MaxLimit:       200,
			AvailableLimit: 200,
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

	accountRepoMock := repository.NewMockAccountRepository(ctrl)

	accountRepoMock.EXPECT().Retrieve(gomock.Any()).Return(&mockAccount, nil)

	as := NewAccount(accountRepoMock)

	_, violations := as.InitAccount(true, 200)

	assert.Equal(t, violations, []string{"account-already-initialized"})
}
