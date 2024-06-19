package handler

import "github.com/posilva/simpleboards/internal/core/domain"

// PutScore ...
type PutScore struct {
	Entry    string          `json:"entry"`
	Score    float64         `json:"score"`
	Metadata domain.Metadata `json:"metadata"`
}
