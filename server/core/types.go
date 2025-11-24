package core

import "context"

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ServiceConfig interface {
	Empty() bool
}

type Service interface {
	Name() string
	NeedsRefresh() bool
	Refresh(context.Context) error
	StateForPrompt() *string
}
