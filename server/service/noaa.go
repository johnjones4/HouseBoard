package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/core"
	"net/http"
	"strings"
	"time"
)

type Weather struct {
	RadarURL string                              `json:"radarURL"`
	Forecast []NOAAWeatherForecastPeriod         `json:"forecast"`
	Alerts   []NOAAWeatherAlertFeatureProperties `json:"alerts"`
}

type NOAAWeatherPointProperties struct {
	ForecastURL  string `json:"forecast"`
	ForecastZone string `json:"forecastZone"`
	RadarStation string `json:"radarStation"`
}

type NOAAWeatherPointResponse struct {
	Properties NOAAWeatherPointProperties `json:"properties"`
}

type NOAAWeatherForecastPeriod struct {
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

type NOAAWeatherForecastProperties struct {
	Periods []NOAAWeatherForecastPeriod `json:"periods"`
}

type NOAAWeatherForecastResponse struct {
	Properties NOAAWeatherForecastProperties `json:"properties"`
}

type NOAAWeatherAlertFeatureProperties struct {
	ID            string   `json:"id"`
	AffectedZones []string `json:"affectedZones"`
	Headline      string   `json:"headline"`
}

type NOAAWeatherAlertFeature struct {
	Properties NOAAWeatherAlertFeatureProperties `json:"properties"`
}

type NOAAWeatherAlertResponse struct {
	Features []NOAAWeatherAlertFeature `json:"features"`
}

type NOAAConfiguration struct {
	Location core.Coordinate
}

func (c NOAAConfiguration) Empty() bool {
	return c.Location.Latitude == 0 && c.Location.Longitude == 0
}

func (c NOAAConfiguration) Service() core.Service {
	return &noaa{c}
}

type noaa struct {
	configuration NOAAConfiguration
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

func (f *noaa) predictWeather(coord core.Coordinate) (Weather, error) {
	point, err := makeWeatherAPIPointRequest(coord)
	if err != nil {
		return Weather{}, err
	}

	radarURL := fmt.Sprintf("https://radar.weather.gov/ridge/lite/%s_loop.gif?v=%d", point.RadarStation, time.Now().Unix())

	forecast, err := makeWeatherAPIForecastCall(point)
	if err != nil {
		return Weather{}, err
	}

	alerts, err := makeWeatherAPIAlertCall(point)
	if err != nil {
		return Weather{}, err
	}

	return Weather{
		RadarURL: radarURL,
		Forecast: forecast,
		Alerts:   alerts,
	}, nil
}

func makeWeatherAPIPointRequest(coord core.Coordinate) (NOAAWeatherPointProperties, error) {
	httpResponse, err := http.Get(fmt.Sprintf("https://api.weather.gov/points/%f,%f", coord.Latitude, coord.Longitude))
	if err != nil {
		return NOAAWeatherPointProperties{}, err
	}

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return NOAAWeatherPointProperties{}, err
	}

	var pointResponse NOAAWeatherPointResponse
	err = json.Unmarshal(responseBytes, &pointResponse)

	return pointResponse.Properties, err
}

func makeWeatherAPIForecastCall(point NOAAWeatherPointProperties) ([]NOAAWeatherForecastPeriod, error) {
	httpResponse, err := http.Get(point.ForecastURL)
	if err != nil {
		return nil, nil
	}

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, nil
	}

	var response NOAAWeatherForecastResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return nil, nil
	}

	return response.Properties.Periods, nil
}

func makeWeatherAPIAlertCall(point NOAAWeatherPointProperties) ([]NOAAWeatherAlertFeatureProperties, error) {
	zoneId := strings.Replace(point.ForecastZone, "https://api.weather.gov/zones/forecast/", "", 1)

	httpResponse, err := http.Get(fmt.Sprintf("https://api.weather.gov/alerts/active/zone/%s", zoneId))
	if err != nil {
		return nil, err
	}

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	var response NOAAWeatherAlertResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return nil, err
	}

	featureProps := make([]NOAAWeatherAlertFeatureProperties, 0)
	for _, feature := range response.Features {
		for _, zone := range feature.Properties.AffectedZones {
			if zone == point.ForecastZone {
				featureProps = append(featureProps, feature.Properties)
			}
		}
	}
	return featureProps, nil
}
