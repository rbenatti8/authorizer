package ports

import (
	"github.com/authorizer/internal/core/domain"
	"time"
)

type AccountRepository interface {
	Create(account domain.Account) error
	Retrieve(currentTime time.Time) (*domain.Account, error)
	Update(account domain.Account) error
}
