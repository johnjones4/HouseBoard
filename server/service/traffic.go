package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const maxMonthlyQueries = 20000

type TrafficDestination struct {
	Destination       string `json:"destination"`
	ExpectedDuration  int    `json:"expectedDuration"`
	EstimatedDuration int    `json:"estimatedDuration"`
}

type TrafficResp struct {
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

func (c trafficConfiguration) Service() *Traffic {
	return &Traffic{
		configuration: c,
		delay:         time.Second * time.Duration((60*60*24*31)/(maxMonthlyQueries/len(c.Destinations))),
	}
}

type Traffic struct {
	configuration trafficConfiguration
	lastLoad      time.Time
	delay         time.Duration
	TrafficResp   *TrafficResp
}

func (t *Traffic) Name() string {
	return "traffic"
}

func (t *Traffic) Refresh(c context.Context) error {
	tr := TrafficResp{
		Destinations: make([]TrafficDestination, len(t.configuration.Destinations)),
	}
	for i, destination := range t.configuration.Destinations {
		params := make(url.Values)
		params.Set("origin", t.configuration.Origin)
		params.Set("destination", destination.Destination)
		params.Set("key", t.configuration.APIKey)

		res, err := http.Get("https://maps.googleapis.com/maps/api/directions/json?" + params.Encode())
		if err != nil {
			return err
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		var responseBody trafficResponse
		err = json.Unmarshal(body, &responseBody)
		if err != nil {
			return err
		}

		if len(responseBody.Routes) == 0 || len(responseBody.Routes[0].Legs) == 0 {
			return errors.New("bad response from directions service")
		}

		tr.Destinations[i] = TrafficDestination{
			Destination:       destination.Destination,
			ExpectedDuration:  destination.ExpectedDuration,
			EstimatedDuration: responseBody.Routes[0].Legs[0].Duration.Value,
		}
	}
	t.lastLoad = time.Now()
	t.TrafficResp = &tr
	return nil
}

func (t *Traffic) NeedsRefresh() bool {
	return time.Now().After(t.lastLoad.Add(t.delay))
}

func (t *Traffic) StateForPrompt() *string {
	if t.TrafficResp == nil {
		return nil
	}

	var summary strings.Builder

	summary.WriteString("Traffic Conditions:\n")

	for _, item := range t.TrafficResp.Destinations {
		summary.WriteString(fmt.Sprintf("- %s Travel Time: %s (Expected: %s)\n", item.Destination, (time.Duration(item.EstimatedDuration) * time.Second).String(), (time.Duration(item.ExpectedDuration) * time.Second).String()))
	}

	str := summary.String()

	return &str
}
