package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"main/core"
	"net/http"
	"net/url"
	"time"
)

const maxMonthlyQueries = 20000

type TrafficDestination struct {
	Destination       string `json:"destination"`
	ExpectedDuration  int    `json:"expectedDuration"`
	EstimatedDuration int    `json:"estimatedDuration"`
}

type Traffic struct {
	Destinations []TrafficDestination `json:"destinations"`
}

type trafficLegDurationResponse struct {
	Value int `json:"value"`
}

type trafficLegResponse struct {
	Duration trafficLegDurationResponse `json:"duration"`
}

type trafficRouteResponse struct {
	Legs []trafficLegResponse `json:"legs"`
}

type trafficResponse struct {
	Routes []trafficRouteResponse `json:"routes"`
}

type trafficDestination struct {
	Destination      string `json:"destination"`
	ExpectedDuration int    `json:"expectedDuration"`
}

type trafficConfiguration struct {
	APIKey       string               `json:"apiKey"`
	Origin       string               `json:"origin"`
	Destinations []trafficDestination `json:"destinations"`
}

func (c trafficConfiguration) Empty() bool {
	return c.Origin == "" && len(c.Destinations) == 0
}

func (c trafficConfiguration) Service() core.Service {
	return &traffic{
		configuration: c,
		delay:         time.Second * time.Duration((60*60*24*31)/(maxMonthlyQueries/len(c.Destinations))),
	}
}

type traffic struct {
	configuration trafficConfiguration
	lastLoad      time.Time
	delay         time.Duration
}

func (t *traffic) Name() string {
	return "traffic"
}

func (t *traffic) Info(c context.Context) (interface{}, error) {
	tr := Traffic{
		Destinations: make([]TrafficDestination, len(t.configuration.Destinations)),
	}
	for i, destination := range t.configuration.Destinations {
		params := make(url.Values)
		params.Set("origin", t.configuration.Origin)
		params.Set("destination", destination.Destination)
		params.Set("key", t.configuration.APIKey)

		res, err := http.Get("https://maps.googleapis.com/maps/api/directions/json?" + params.Encode())
		if err != nil {
			return nil, err
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var responseBody trafficResponse
		err = json.Unmarshal(body, &responseBody)
		if err != nil {
			return nil, err
		}

		if len(responseBody.Routes) == 0 || len(responseBody.Routes[0].Legs) == 0 {
			return nil, errors.New("bad response from directions service")
		}

		tr.Destinations[i] = TrafficDestination{
			Destination:       destination.Destination,
			ExpectedDuration:  destination.ExpectedDuration,
			EstimatedDuration: responseBody.Routes[0].Legs[0].Duration.Value,
		}
	}
	t.lastLoad = time.Now()
	return tr, nil
}

func (t *traffic) NeedsRefresh() bool {
	return time.Now().After(t.lastLoad.Add(t.delay))
}
