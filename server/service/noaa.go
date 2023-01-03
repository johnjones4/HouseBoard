package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"main/core"
	"net/http"
	"strings"
	"time"
)

type weather struct {
	RadarURL string                              `json:"radarURL"`
	Forecast []noaaWeatherForecastPeriod         `json:"forecast"`
	Alerts   []noaaWeatherAlertFeatureProperties `json:"alerts"`
}

type noaaWeatherPointProperties struct {
	ForecastURL  string `json:"forecast"`
	ForecastZone string `json:"forecastZone"`
	RadarStation string `json:"radarStation"`
}

type noaaWeatherPointResponse struct {
	Properties noaaWeatherPointProperties `json:"properties"`
}

type noaaWeatherForecastPeriod struct {
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	DetailedForecast string    `json:"detailedForecast"`
	Name             string    `json:"name"`
	Temperature      float64   `json:"temperature"`
	TemperatureUnit  string    `json:"temperatureUnit"`
	WindSpeed        string    `json:"windSpeed"`
	WindDirection    string    `json:"windDirection"`
	Icon             string    `json:"icon"`
	IsDaytime        bool      `json:"isDaytime"`
}

type noaaWeatherForecastProperties struct {
	Periods []noaaWeatherForecastPeriod `json:"periods"`
}

type noaaWeatherForecastResponse struct {
	Properties noaaWeatherForecastProperties `json:"properties"`
}

type noaaWeatherAlertFeatureProperties struct {
	ID            string   `json:"id"`
	AffectedZones []string `json:"affectedZones"`
	Headline      string   `json:"headline"`
}

type noaaWeatherAlertFeature struct {
	Properties noaaWeatherAlertFeatureProperties `json:"properties"`
}

type noaaWeatherAlertResponse struct {
	Features []noaaWeatherAlertFeature `json:"features"`
}

type noaaConfiguration struct {
	Location core.Coordinate
}

func (c noaaConfiguration) Empty() bool {
	return c.Location.Latitude == 0 && c.Location.Longitude == 0
}

func (c noaaConfiguration) Service() core.Service {
	return &noaa{c}
}

type noaa struct {
	configuration noaaConfiguration
}

func (f *noaa) Name() string {
	return "noaa"
}

func (f *noaa) Info(c context.Context) (interface{}, error) {
	weather, err := f.predictWeather(f.configuration.Location)
	if err != nil {
		return nil, err
	}
	return weather, nil
}

func (f *noaa) NeedsRefresh() bool {
	return true
}

func (f *noaa) predictWeather(coord core.Coordinate) (weather, error) {
	point, err := makeWeatherAPIPointRequest(coord)
	if err != nil {
		return weather{}, err
	}

	radarURL := fmt.Sprintf("https://radar.weather.gov/ridge/standard/%s_loop.gif?refreshed=%d", point.RadarStation, time.Now().Unix())

	forecast, err := makeWeatherAPIForecastCall(point)
	if err != nil {
		return weather{}, err
	}

	alerts, err := makeWeatherAPIAlertCall(point)
	if err != nil {
		return weather{}, err
	}

	return weather{
		RadarURL: radarURL,
		Forecast: forecast,
		Alerts:   alerts,
	}, nil
}

func makeWeatherAPIPointRequest(coord core.Coordinate) (noaaWeatherPointProperties, error) {
	httpResponse, err := http.Get(fmt.Sprintf("https://api.weather.gov/points/%f,%f", coord.Latitude, coord.Longitude))
	if err != nil {
		return noaaWeatherPointProperties{}, err
	}

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return noaaWeatherPointProperties{}, err
	}

	var pointResponse noaaWeatherPointResponse
	err = json.Unmarshal(responseBytes, &pointResponse)

	return pointResponse.Properties, err
}

func makeWeatherAPIForecastCall(point noaaWeatherPointProperties) ([]noaaWeatherForecastPeriod, error) {
	httpResponse, err := http.Get(point.ForecastURL)
	if err != nil {
		return nil, err
	}

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	var response noaaWeatherForecastResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Properties.Periods) == 0 {
		return nil, errors.New("no forecast data returned")
	}

	return response.Properties.Periods, nil
}

func makeWeatherAPIAlertCall(point noaaWeatherPointProperties) ([]noaaWeatherAlertFeatureProperties, error) {
	zoneId := strings.Replace(point.ForecastZone, "https://api.weather.gov/zones/forecast/", "", 1)

	httpResponse, err := http.Get(fmt.Sprintf("https://api.weather.gov/alerts/active/zone/%s", zoneId))
	if err != nil {
		return nil, err
	}

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	var response noaaWeatherAlertResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return nil, err
	}

	featureProps := make([]noaaWeatherAlertFeatureProperties, 0)
	for _, feature := range response.Features {
		for _, zone := range feature.Properties.AffectedZones {
			if zone == point.ForecastZone {
				featureProps = append(featureProps, feature.Properties)
			}
		}
	}
	return featureProps, nil
}
