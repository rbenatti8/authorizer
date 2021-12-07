package dto

type AccountOperation struct {
	ActiveCard     bool  `json:"active-card"`
	AvailableLimit int64 `json:"available-limit"`
}
