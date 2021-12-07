package cli

import (
	"encoding/json"
	"fmt"
	"github.com/authorizer/internal/core/domain"
	"github.com/authorizer/internal/core/service"
	"github.com/authorizer/internal/dto"
	"io"
)

type Handler struct {
	accountService     service.Account
	transactionService service.Transaction
}

func NewHandler(as service.Account, ts service.Transaction) Handler {
	return Handler{accountService: as, transactionService: ts}
}

func (h Handler) Handle(r io.Reader, w io.Writer) error {
	decoder := json.NewDecoder(r)
	decoder.UseNumber()

	for {
		var input dto.Input

		err := decoder.Decode(&input)

		if err == io.EOF {
			break
		}

		resp := h.handle(input)

		jm, _ := json.Marshal(resp)

		_, _ = fmt.Fprintln(w, string(jm))
	}

	return nil
}

func (h Handler) handle(input dto.Input) dto.Output {
	var (
		account    *domain.Account
		violations []string
	)

	if input.Account != nil {
		account, violations = h.accountService.InitAccount(input.Account.ActiveCard, input.Account.AvailableLimit)
		return buildOutput(account, violations)
	}

	if input.Transaction != nil {
		t := domain.Transaction{
			Merchant: input.Transaction.Merchant,
			Amount:   input.Transaction.Amount,
			Time:     input.Transaction.Time,
		}

		account, violations = h.transactionService.Authorize(t)
		return buildOutput(account, violations)
	}

	return dto.Output{}
}

func buildOutput(account *domain.Account, violations []string) dto.Output {
	output := dto.Output{
		Violations: violations,
	}

	if account != nil {
		output.Account = dto.AccountOutput{
			ActiveCard:     &account.Ledger.ActiveCard,
			AvailableLimit: &account.Ledger.AvailableLimit,
		}
	}

	return output
}
