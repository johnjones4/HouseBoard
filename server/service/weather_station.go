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
	WeatherStatonResponse *WeatherStatonResponse
}

type WeatherStatonResponse struct {
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
	Items []WeatherStatonResponse `json:"items"`
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

	var info weatherStationResponseBody
	err = json.Unmarshal(body, &info)
	if err != nil {
		return err
	}

	if len(info.Items) == 0 {
		return nil
	}

	w.WeatherStatonResponse = &info.Items[len(info.Items)-1]

	return nil
}

func (w *WeatherStation) NeedsRefresh() bool {
	return true
}

func (w *WeatherStation) StateForPrompt() *string {
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
		w.WeatherStatonResponse.AvgWindSpeed,
		w.WeatherStatonResponse.VaneDirection,
		w.WeatherStatonResponse.RelativeHumidity,
		w.WeatherStatonResponse.Pressure,
	)
	return &str
}
