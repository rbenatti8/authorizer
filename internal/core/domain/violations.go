package domain

var (
	AccountNotInitializedViolation     = "account-not-initialized"
	CardNotActiveViolation             = "card-not-active"
	InsufficientLimitViolation         = "insufficient-limit"
	AccountAlreadyInitializedViolation = "account-already-initialized"
	DoubledTransactionViolation        = "doubled-transaction"
)

type Violations []string

func (v *Violations) AddViolation(violations ...string) {
	if len(violations) > 0 {
		v.add(violations)
	}
}

func (v *Violations) add(violations []string) {
	for _, violation := range violations {
		if violation != "" {
			*v = append(*v, violation)
		}
	}
}
