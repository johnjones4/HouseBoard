package service

import (
	"context"
	"encoding/json"
	"io"
	"main/core"
	"net/http"
	"time"
)

type weatherStationConfiguration struct {
	Upstream string `json:"upstream"`
}

func (c weatherStationConfiguration) Empty() bool {
	return c.Upstream == ""
}

func (c weatherStationConfiguration) Service() core.Service {
	return &weatherStation{
		configuration: c,
	}
}

type weatherStation struct {
	configuration weatherStationConfiguration
}

type weatherStatonResponse struct {
	Timestamp        time.Time `json:"timestamp"`
	Uptime           int64     `json:"uptime"`
	AvgWindSpeed     float64   `json:"avg_wind_speed"`
	MinWindSpeed     float64   `json:"min_wind_speed"`
	MaxWindSpeed     float64   `json:"max_wind_speed"`
	Temperature      float64   `json:"temperature"`
	Gas              float64   `json:"gas"`
	RelativeHumidity float64   `json:"relative_humidity"`
	Pressure         float64   `json:"pressure"`
}

func (w *weatherStation) Name() string {
	return "weatherstation"
}

func (w *weatherStation) Info(c context.Context) (interface{}, error) {
	res, err := http.Get(w.configuration.Upstream)
	if err != nil {
		return weatherStatonResponse{}, nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return weatherStatonResponse{}, nil
	}

	var info weatherStatonResponse
	err = json.Unmarshal(body, &info)
	if err != nil {
		return weatherStatonResponse{}, nil
	}

	return info, nil
}

func (w *weatherStation) NeedsRefresh() bool {
	return true
}
