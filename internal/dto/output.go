package dto

type AccountOutput struct {
	ActiveCard     *bool  `json:"active-card,omitempty"`
	AvailableLimit *int64 `json:"available-limit,omitempty"`
}

type Output struct {
	Account    AccountOutput `json:"account"`
	Violations []string      `json:"violations"`
}
