package handler

type PutScore struct {
	Entry string  `json:"entry"`
	Score float64 `json:"score"`
}
