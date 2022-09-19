package core

import "context"

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ServiceConfig interface {
	Empty() bool
	Service() Service
}

type Service interface {
	Name() string
	Info(c context.Context) (interface{}, error)
	NeedsRefresh() bool
}
