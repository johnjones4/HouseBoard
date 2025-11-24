package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type weatherStationConfiguration struct {
	Upstream string `json:"upstream"`
}

func (c weatherStationConfiguration) Empty() bool {
	return c.Upstream == ""
}

func (c weatherStationConfiguration) Service() *WeatherStation {
	return &WeatherStation{
		configuration: c,
	}
}

type WeatherStation struct {
	configuration         weatherStationConfiguration
	WeatherStatonResponse *WeatherReading
	LastUpdated           *time.Time
}

type weatherAverage struct {
	Period        string         `json:"period"`
	Start         time.Time      `json:"start"`
	Averages      WeatherReading `json:"averages"`
	RainfallTotal float64        `json:"rainfallTotal"`
}

type weatherResponse struct {
	Readings []WeatherItem    `json:"readings"`
	Periods  []weatherAverage `json:"periods"`
}

type WeatherReading struct {
	WindSpeed     *float64 `json:"windSpeed"`
	VaneDirection *float64 `json:"vaneDirection"`
	Temperature   *float64 `json:"temperature"`
	Pressure      *float64 `json:"pressure"`
	Humidity      *float64 `json:"humidity"`
	Gas           *float64 `json:"gas"`
	Rainfall      *float64 `json:"rainfall"`
}

type WeatherPayload struct {
	Source string `json:"source"`
	WeatherReading
}

type WeatherItem struct {
	Timestamp time.Time `json:"timestamp"`
	WeatherPayload
}

func (w *WeatherStation) Name() string {
	return "weatherstation"
}

func (w *WeatherStation) Refresh(c context.Context) error {
	res, err := http.Get(w.configuration.Upstream)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var info weatherResponse
	err = json.Unmarshal(body, &info)
	if err != nil {
		return err
	}

	if len(info.Readings) == 0 {
		return nil
	}

	w.LastUpdated = &info.Readings[len(info.Readings)-1].Timestamp

	w.WeatherStatonResponse = &WeatherReading{}

	for _, p := range info.Periods {
		if p.Period == "5min" {
			w.WeatherStatonResponse.Gas = p.Averages.Gas
			w.WeatherStatonResponse.Temperature = p.Averages.Temperature
			w.WeatherStatonResponse.Humidity = p.Averages.Humidity
			w.WeatherStatonResponse.WindSpeed = p.Averages.WindSpeed
			w.WeatherStatonResponse.VaneDirection = p.Averages.VaneDirection
			w.WeatherStatonResponse.Pressure = p.Averages.Pressure
		}
		if p.Period == "1hr" {
			w.WeatherStatonResponse.Rainfall = &p.RainfallTotal
		}
	}

	return nil
}

func (w *WeatherStation) NeedsRefresh() bool {
	return w != nil
}

func (w *WeatherStation) StateForPrompt() *string {
	if w == nil {
		return nil
	}

	if w.WeatherStatonResponse == nil {
		return nil
	}

	str := fmt.Sprintf(`Household Weather Station Reading:
- Temperature: %0.2f degrees Fahrenheit
- Wind Speed %0.2f MPH
- Wind Direction (Degrees): %0.2f
- Relative Humidity: %0.2f%%
- Atmospheric Pressure: %0.2f inHg`,
		w.WeatherStatonResponse.Temperature,
		w.WeatherStatonResponse.WindSpeed,
		w.WeatherStatonResponse.VaneDirection,
		w.WeatherStatonResponse.Humidity,
		w.WeatherStatonResponse.Pressure,
	)
	return &str
}
