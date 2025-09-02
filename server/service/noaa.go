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

type Weather struct {
	LocalRadarURL    string `json:"radarURL"`
	NationalRadarURL string
	Forecast         []NoaaWeatherForecastPeriod         `json:"forecast"`
	Alerts           []NoaaWeatherAlertFeatureProperties `json:"alerts"`
}

type NoaaWeatherPointProperties struct {
	ForecastURL  string `json:"forecast"`
	ForecastZone string `json:"forecastZone"`
	RadarStation string `json:"radarStation"`
}

type noaaWeatherPointResponse struct {
	Properties NoaaWeatherPointProperties `json:"properties"`
}

type NoaaWeatherValue struct {
	Value    float64 `json:"value"`
	UnitCode string  `json:"unitCode"`
}

type NoaaWeatherForecastPeriod struct {
	StartTime                  time.Time        `json:"startTime"`
	EndTime                    time.Time        `json:"endTime"`
	DetailedForecast           string           `json:"detailedForecast"`
	Name                       string           `json:"name"`
	Temperature                float64          `json:"temperature"`
	TemperatureUnit            string           `json:"temperatureUnit"`
	WindSpeed                  string           `json:"windSpeed"`
	WindDirection              string           `json:"windDirection"`
	Icon                       string           `json:"icon"`
	IsDaytime                  bool             `json:"isDaytime"`
	ProbabilityOfPrecipitation NoaaWeatherValue `json:"probabilityOfPrecipitation"`
	RelativeHumidity           NoaaWeatherValue `json:"relativeHumidity"`
	Dewpoint                   NoaaWeatherValue `json:"dewpoint"`
}

type noaaWeatherForecastProperties struct {
	Periods []NoaaWeatherForecastPeriod `json:"periods"`
}

type noaaWeatherForecastResponse struct {
	Properties noaaWeatherForecastProperties `json:"properties"`
}

type NoaaWeatherAlertFeatureProperties struct {
	ID            string   `json:"id"`
	AffectedZones []string `json:"affectedZones"`
	Headline      string   `json:"headline"`
}

type noaaWeatherAlertFeature struct {
	Properties NoaaWeatherAlertFeatureProperties `json:"properties"`
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

func (c noaaConfiguration) Service() *NOAA {
	return &NOAA{c, nil}
}

type NOAA struct {
	configuration noaaConfiguration
	Weather       *Weather
}

func (f *NOAA) Name() string {
	return "noaa"
}

func (f *NOAA) Refresh(c context.Context) error {
	weather, err := f.predictWeather(f.configuration.Location)
	if err != nil {
		return err
	}
	f.Weather = &weather
	return nil
}

func (f *NOAA) NeedsRefresh() bool {
	return f != nil
}

func (f *NOAA) predictWeather(coord core.Coordinate) (Weather, error) {
	point, err := makeWeatherAPIPointRequest(coord)
	if err != nil {
		return Weather{}, err
	}

	forecast, err := makeWeatherAPIForecastCall(point)
	if err != nil {
		return Weather{}, err
	}

	alerts, err := makeWeatherAPIAlertCall(point)
	if err != nil {
		return Weather{}, err
	}

	return Weather{
		LocalRadarURL:    fmt.Sprintf("https://radar.weather.gov/ridge/standard/%s_loop.gif", point.RadarStation),
		NationalRadarURL: "https://radar.weather.gov/ridge/standard/CONUS-LARGE_loop.gif",
		Forecast:         forecast,
		Alerts:           alerts,
	}, nil
}

func makeWeatherAPIPointRequest(coord core.Coordinate) (NoaaWeatherPointProperties, error) {
	httpResponse, err := http.Get(fmt.Sprintf("https://api.weather.gov/points/%f,%f", coord.Latitude, coord.Longitude))
	if err != nil {
		return NoaaWeatherPointProperties{}, err
	}

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return NoaaWeatherPointProperties{}, err
	}

	var pointResponse noaaWeatherPointResponse
	err = json.Unmarshal(responseBytes, &pointResponse)

	return pointResponse.Properties, err
}

func makeWeatherAPIForecastCall(point NoaaWeatherPointProperties) ([]NoaaWeatherForecastPeriod, error) {
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

func makeWeatherAPIAlertCall(point NoaaWeatherPointProperties) ([]NoaaWeatherAlertFeatureProperties, error) {
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

	featureProps := make([]NoaaWeatherAlertFeatureProperties, 0)
	for _, feature := range response.Features {
		for _, zone := range feature.Properties.AffectedZones {
			if zone == point.ForecastZone {
				featureProps = append(featureProps, feature.Properties)
			}
		}
	}
	return featureProps, nil
}

func (n *NOAA) StateForPrompt() *string {
	return nil
}
