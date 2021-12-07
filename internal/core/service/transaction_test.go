package service

import (
	"github.com/authorizer/internal/core/domain"
	"github.com/authorizer/internal/driven/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTransaction_Authorize_Without_Violations(t *testing.T) {
	testCases := []struct {
		name            string
		mockAccount     domain.Account
		transaction     domain.Transaction
		expectedAccount domain.Account
	}{
		{
			name: "Processando uma transação com sucesso",
			transaction: domain.Transaction{
				Merchant: "xablau testador",
				Amount:   100,
				Time:     time.Date(2021, 10, 10, 10, 0, 0, 0, time.Local),
			},
			mockAccount: domain.Account{
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
								Duration:           2 * time.Minute,
								CurrentPeriodUsed:  0,
								CurrentPeriodSpend: 0,
								PeriodEndsDate:     time.Date(2021, 10, 10, 10, 2, 0, 0, time.Local),
							},
							RuleViolation: "high-frequency-small-interval",
						},
					},
				},
				Authorizations: []domain.TransactionAuthorization{},
			},
			expectedAccount: domain.Account{
				Ledger: domain.Ledger{
					ActiveCard:     true,
					MaxLimit:       200,
					AvailableLimit: 100,
				},
				SpendingControl: domain.SpendingControl{
					Rules: []domain.Rule{
						{
							Name:       "max transactions in 2 minutes",
							Type:       "usage-limit",
							UsageLimit: 3,
							Accumulator: &domain.Accumulator{
								Duration:           2 * time.Minute,
								CurrentPeriodUsed:  1,
								CurrentPeriodSpend: 100,
								PeriodEndsDate:     time.Date(2021, 10, 10, 10, 2, 0, 0, time.Local),
							},
							RuleViolation: "high-frequency-small-interval",
						},
					},
				},
				Authorizations: []domain.TransactionAuthorization{
					{
						Merchant:       "xablau testador",
						Amount:         100,
						AvailableLimit: 200,
						Time:           time.Date(2021, 10, 10, 10, 0, 0, 0, time.Local),
					},
				},
			},
		},
		{
			name: "Processando a terceira transação com sucesso",
			transaction: domain.Transaction{
				Merchant: "Merchant3",
				Amount:   25,
				Time:     time.Date(2021, 10, 10, 10, 1, 0, 0, time.Local),
			},
			mockAccount: domain.Account{
				Ledger: domain.Ledger{
					ActiveCard:     true,
					MaxLimit:       500,
					AvailableLimit: 450,
				},
				SpendingControl: domain.SpendingControl{
					Rules: []domain.Rule{
						{
							Name:       "max transactions in 2 minutes",
							Type:       "usage-limit",
							UsageLimit: 3,
							Accumulator: &domain.Accumulator{
								Duration:           2 * time.Minute,
								CurrentPeriodUsed:  2,
								CurrentPeriodSpend: 50,
								PeriodEndsDate:     time.Date(2021, 10, 10, 10, 2, 0, 0, time.Local),
							},
							RuleViolation: "high-frequency-small-interval",
						},
					},
				},
				Authorizations: []domain.TransactionAuthorization{
					{
						Merchant:       "Merchant1",
						Amount:         25,
						AvailableLimit: 500,
						Time:           time.Date(2021, 10, 10, 10, 0, 0, 0, time.Local),
					},
					{
						Merchant:       "Merchant2",
						Amount:         25,
						AvailableLimit: 475,
						Time:           time.Date(2021, 10, 10, 10, 0, 30, 0, time.Local),
					},
				},
			},
			expectedAccount: domain.Account{
				Ledger: domain.Ledger{
					ActiveCard:     true,
					MaxLimit:       500,
					AvailableLimit: 425,
				},
				SpendingControl: domain.SpendingControl{
					Rules: []domain.Rule{
						{
							Name:       "max transactions in 2 minutes",
							Type:       "usage-limit",
							UsageLimit: 3,
							Accumulator: &domain.Accumulator{
								Duration:           2 * time.Minute,
								CurrentPeriodUsed:  3,
								CurrentPeriodSpend: 75,
								PeriodEndsDate:     time.Date(2021, 10, 10, 10, 2, 0, 0, time.Local),
							},
							RuleViolation: "high-frequency-small-interval",
						},
					},
				},
				Authorizations: []domain.TransactionAuthorization{
					{
						Merchant:       "Merchant1",
						Amount:         25,
						AvailableLimit: 500,
						Time:           time.Date(2021, 10, 10, 10, 0, 0, 0, time.Local),
					},
					{
						Merchant:       "Merchant2",
						Amount:         25,
						AvailableLimit: 475,
						Time:           time.Date(2021, 10, 10, 10, 0, 30, 0, time.Local),
					},
					{
						Merchant:       "Merchant3",
						Amount:         25,
						AvailableLimit: 450,
						Time:           time.Date(2021, 10, 10, 10, 1, 0, 0, time.Local),
					},
				},
			},
		},
		{
			name: "Processando transação com intervalo de tempo maior",
			transaction: domain.Transaction{
				Merchant: "Merchant3",
				Amount:   25,
				Time:     time.Date(2021, 10, 10, 10, 10, 0, 0, time.Local),
			},
			mockAccount: domain.Account{
				Ledger: domain.Ledger{
					ActiveCard:     true,
					MaxLimit:       500,
					AvailableLimit: 450,
				},
				SpendingControl: domain.SpendingControl{
					Rules: []domain.Rule{
						{
							Name:       "max transactions in 2 minutes",
							Type:       "usage-limit",
							UsageLimit: 3,
							Accumulator: &domain.Accumulator{
								Duration:           2 * time.Minute,
								CurrentPeriodUsed:  2,
								CurrentPeriodSpend: 50,
								PeriodEndsDate:     time.Date(2021, 10, 10, 10, 2, 0, 0, time.Local),
							},
							RuleViolation: "high-frequency-small-interval",
						},
					},
				},
				Authorizations: []domain.TransactionAuthorization{
					{
						Merchant:       "Merchant1",
						Amount:         25,
						AvailableLimit: 500,
						Time:           time.Date(2021, 10, 10, 10, 0, 0, 0, time.Local),
					},
					{
						Merchant:       "Merchant2",
						Amount:         25,
						AvailableLimit: 475,
						Time:           time.Date(2021, 10, 10, 10, 0, 30, 0, time.Local),
					},
				},
			},
			expectedAccount: domain.Account{
				Ledger: domain.Ledger{
					ActiveCard:     true,
					MaxLimit:       500,
					AvailableLimit: 425,
				},
				SpendingControl: domain.SpendingControl{
					Rules: []domain.Rule{
						{
							Name:       "max transactions in 2 minutes",
							Type:       "usage-limit",
							UsageLimit: 3,
							Accumulator: &domain.Accumulator{
								Duration:           2 * time.Minute,
								CurrentPeriodUsed:  3,
								CurrentPeriodSpend: 75,
								PeriodEndsDate:     time.Date(2021, 10, 10, 10, 2, 0, 0, time.Local),
							},
							RuleViolation: "high-frequency-small-interval",
						},
					},
				},
				Authorizations: []domain.TransactionAuthorization{
					{
						Merchant:       "Merchant1",
						Amount:         25,
						AvailableLimit: 500,
						Time:           time.Date(2021, 10, 10, 10, 0, 0, 0, time.Local),
					},
					{
						Merchant:       "Merchant2",
						Amount:         25,
						AvailableLimit: 475,
						Time:           time.Date(2021, 10, 10, 10, 0, 30, 0, time.Local),
					},
					{
						Merchant:       "Merchant3",
						Amount:         25,
						AvailableLimit: 450,
						Time:           time.Date(2021, 10, 10, 10, 10, 0, 0, time.Local),
					},
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			accountRepoMock := repository.NewMockAccountRepository(ctrl)

			accountRepoMock.EXPECT().Retrieve(tt.transaction.Time).Return(&tt.mockAccount, nil)
			accountRepoMock.EXPECT().Update(tt.expectedAccount).Return(nil)

			ts := NewTransaction(accountRepoMock)

			account, _ := ts.Authorize(tt.transaction)

			assert.Equal(t, account.Ledger.AvailableLimit, tt.expectedAccount.Ledger.AvailableLimit)
		})
	}
}

func TestTransaction_Authorize_With_Violations(t *testing.T) {
	testCases := []struct {
		name               string
		mockAccount        *domain.Account
		transaction        domain.Transaction
		expectedViolations domain.Violations
	}{
		{
			name: "Processando uma transação que viola a lógica account-not-initialized",
			transaction: domain.Transaction{
				Merchant: "xablau testador",
				Amount:   100,
				Time:     time.Date(2021, 10, 10, 10, 0, 0, 0, time.Local),
			},
			mockAccount:        nil,
			expectedViolations: domain.Violations{"account-not-initialized"},
		},
		{
			name: "Processando uma transação que viola a lógica card-not-active",
			transaction: domain.Transaction{
				Merchant: "xablau testador",
				Amount:   100,
				Time:     time.Date(2021, 10, 10, 10, 0, 0, 0, time.Local),
			},
			mockAccount: &domain.Account{
				Ledger: domain.Ledger{
					ActiveCard:     false,
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
								Duration:           2 * time.Minute,
								CurrentPeriodUsed:  0,
								CurrentPeriodSpend: 0,
								PeriodEndsDate:     time.Date(2021, 10, 10, 10, 2, 0, 0, time.Local),
							},
							RuleViolation: "high-frequency-small-interval",
						},
					},
				},
				Authorizations: []domain.TransactionAuthorization{},
			},
			expectedViolations: domain.Violations{"card-not-active"},
		},
		{
			name: "Processando uma transação que viola a lógica insufficient-limit",
			transaction: domain.Transaction{
				Merchant: "xablau testador",
				Amount:   500,
				Time:     time.Date(2021, 10, 10, 10, 0, 0, 0, time.Local),
			},
			mockAccount: &domain.Account{
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
								Duration:           2 * time.Minute,
								CurrentPeriodUsed:  0,
								CurrentPeriodSpend: 0,
								PeriodEndsDate:     time.Date(2021, 10, 10, 10, 2, 0, 0, time.Local),
							},
							RuleViolation: "high-frequency-small-interval",
						},
					},
				},
				Authorizations: []domain.TransactionAuthorization{},
			},
			expectedViolations: domain.Violations{"insufficient-limit"},
		},
		{
			name: "Processando uma transação que viola a lógica high-frequency-small-interval",
			transaction: domain.Transaction{
				Merchant: "xablau testador",
				Amount:   25,
				Time:     time.Date(2021, 10, 10, 10, 1, 30, 0, time.Local),
			},
			mockAccount: &domain.Account{
				Ledger: domain.Ledger{
					ActiveCard:     true,
					MaxLimit:       225,
					AvailableLimit: 150,
				},
				SpendingControl: domain.SpendingControl{
					Rules: []domain.Rule{
						{
							Name:       "max transactions in 2 minutes",
							Type:       "usage-limit",
							UsageLimit: 3,
							Accumulator: &domain.Accumulator{
								Duration:           2 * time.Minute,
								CurrentPeriodUsed:  3,
								CurrentPeriodSpend: 75,
								PeriodEndsDate:     time.Date(2021, 10, 10, 10, 2, 0, 0, time.Local),
							},
							RuleViolation: "high-frequency-small-interval",
						},
					},
				},
				Authorizations: []domain.TransactionAuthorization{
					{
						Merchant:       "Merchant1",
						Amount:         25,
						AvailableLimit: 225,
						Time:           time.Date(2021, 10, 10, 10, 0, 0, 0, time.Local),
					},
					{
						Merchant:       "Merchant2",
						Amount:         25,
						AvailableLimit: 200,
						Time:           time.Date(2021, 10, 10, 10, 0, 30, 0, time.Local),
					},
					{
						Merchant:       "Merchant3",
						Amount:         25,
						AvailableLimit: 175,
						Time:           time.Date(2021, 10, 10, 10, 1, 0, 0, time.Local),
					},
				},
			},
			expectedViolations: domain.Violations{"high-frequency-small-interval"},
		},
		{
			name: "Processando uma transação que viola a lógica doubled-transaction",
			transaction: domain.Transaction{
				Merchant: "Merchant2",
				Amount:   25,
				Time:     time.Date(2021, 10, 10, 10, 1, 30, 0, time.Local),
			},
			mockAccount: &domain.Account{
				Ledger: domain.Ledger{
					ActiveCard:     true,
					MaxLimit:       225,
					AvailableLimit: 175,
				},
				SpendingControl: domain.SpendingControl{
					Rules: []domain.Rule{
						{
							Name:       "max transactions in 2 minutes",
							Type:       "usage-limit",
							UsageLimit: 3,
							Accumulator: &domain.Accumulator{
								Duration:           2 * time.Minute,
								CurrentPeriodUsed:  2,
								CurrentPeriodSpend: 50,
								PeriodEndsDate:     time.Date(2021, 10, 10, 10, 2, 0, 0, time.Local),
							},
							RuleViolation: "high-frequency-small-interval",
						},
					},
				},
				Authorizations: []domain.TransactionAuthorization{
					{
						Merchant:       "Merchant1",
						Amount:         25,
						AvailableLimit: 225,
						Time:           time.Date(2021, 10, 10, 10, 0, 0, 0, time.Local),
					},
					{
						Merchant:       "Merchant2",
						Amount:         25,
						AvailableLimit: 200,
						Time:           time.Date(2021, 10, 10, 10, 0, 30, 0, time.Local),
					},
				},
			},
			expectedViolations: domain.Violations{"doubled-transaction"},
		},
		{
			name: "Processando transações que violam multiplas lógicas",
			transaction: domain.Transaction{
				Merchant: "Merchant3",
				Amount:   100,
				Time:     time.Date(2021, 10, 10, 10, 1, 30, 0, time.Local),
			},
			mockAccount: &domain.Account{
				Ledger: domain.Ledger{
					ActiveCard:     true,
					MaxLimit:       225,
					AvailableLimit: 75,
				},
				SpendingControl: domain.SpendingControl{
					Rules: []domain.Rule{
						{
							Name:       "max transactions in 2 minutes",
							Type:       "usage-limit",
							UsageLimit: 3,
							Accumulator: &domain.Accumulator{
								Duration:           2 * time.Minute,
								CurrentPeriodUsed:  3,
								CurrentPeriodSpend: 150,
								PeriodEndsDate:     time.Date(2021, 10, 10, 10, 2, 0, 0, time.Local),
							},
							RuleViolation: "high-frequency-small-interval",
						},
					},
				},
				Authorizations: []domain.TransactionAuthorization{
					{
						Merchant:       "Merchant1",
						Amount:         25,
						AvailableLimit: 225,
						Time:           time.Date(2021, 10, 10, 10, 0, 0, 0, time.Local),
					},
					{
						Merchant:       "Merchant2",
						Amount:         25,
						AvailableLimit: 200,
						Time:           time.Date(2021, 10, 10, 10, 0, 30, 0, time.Local),
					},
					{
						Merchant:       "Merchant3",
						Amount:         100,
						AvailableLimit: 200,
						Time:           time.Date(2021, 10, 10, 10, 1, 0, 0, time.Local),
					},
				},
			},
			expectedViolations: domain.Violations{"insufficient-limit", "high-frequency-small-interval", "doubled-transaction"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			accountRepoMock := repository.NewMockAccountRepository(ctrl)

			accountRepoMock.EXPECT().Retrieve(gomock.Any()).Return(tt.mockAccount, nil)

			ts := NewTransaction(accountRepoMock)

			_, violations := ts.Authorize(tt.transaction)

			assert.Equal(t, tt.expectedViolations, violations)
		})
	}
}
