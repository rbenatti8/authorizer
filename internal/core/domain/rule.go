package domain

type Rule struct {
	Name          string
	Type          string
	UsageLimit    int64
	Accumulator   *Accumulator
	RuleViolation string
}
