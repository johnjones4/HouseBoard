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
	AvgWindSpeed     float64   `json:"anemometerAverage"`
	MinWindSpeed     float64   `json:"anemometerMin"`
	MaxWindSpeed     float64   `json:"anemometerMax"`
	Temperature      float64   `json:"temperature"`
	Gas              float64   `json:"gas"`
	RelativeHumidity float64   `json:"relativeHumidity"`
	Pressure         float64   `json:"pressure"`
	VaneDirection    float64   `json:"vaneDirection"`
}

type weatherStationResponseBody struct {
	Items []weatherStatonResponse `json:"items"`
}

func (w *weatherStation) Name() string {
	return "weatherstation"
}

func (w *weatherStation) Info(c context.Context) (interface{}, error) {
	res, err := http.Get(w.configuration.Upstream)
	if err != nil {
		return weatherStatonResponse{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return weatherStatonResponse{}, err
	}

	var info weatherStationResponseBody
	err = json.Unmarshal(body, &info)
	if err != nil {
		return weatherStatonResponse{}, err
	}

	if len(info.Items) == 0 {
		return weatherStatonResponse{}, nil
	}

	return info.Items[len(info.Items)-1], nil
}

func (w *weatherStation) NeedsRefresh() bool {
	return true
}
