package service

import (
	"context"
	"encoding/json"
	"fmt"
	"main/core"
	"net/http"
	"net/url"
	"time"
)

type owmResponse struct {
	List []owmResponseItem `json:"list"`
}

type owmResponseItem struct {
	Dt         int64        `json:"dt"`
	Main       owmMain      `json:"main"`
	Weather    []owmWeather `json:"weather"`
	Clouds     owmClouds    `json:"clouds"`
	Wind       owmWind      `json:"wind"`
	Visibility int          `json:"visibility"`
	Pop        float64      `json:"pop"`
	Sys        owmSys       `json:"sys"`
	DtTxt      string       `json:"dt_txt"`
}

// Main represents the main temperature details.
type owmMain struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
	Humidity  int     `json:"humidity"`
	TempKf    float64 `json:"temp_kf"`
}

// Weather represents the weather condition details.
type owmWeather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// Clouds represents the cloud coverage details.
type owmClouds struct {
	All int `json:"all"`
}

// Wind represents the wind details.
type owmWind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

// Sys represents the system details.
type owmSys struct {
	Pod string `json:"pod"`
}

type owmConfiguration struct {
	Location core.Coordinate
	ApiKey   string `json:"apiKey"`
}

func (c owmConfiguration) Empty() bool {
	return c.Location.Latitude == 0 && c.Location.Longitude == 0
}

func (c owmConfiguration) Service() core.Service {
	return &openWeatherMap{
		owmConfiguration: c,
	}
}

const pullInterval = time.Hour

type openWeatherMap struct {
	owmConfiguration
	lastPull time.Time
}

func (f *openWeatherMap) Name() string {
	return "openWeatherMap"
}

func (f *openWeatherMap) Info(c context.Context) (interface{}, error) {
	qs := url.Values{
		"lat":   []string{fmt.Sprint(f.Location.Latitude)},
		"lon":   []string{fmt.Sprint(f.Location.Longitude)},
		"appid": []string{f.ApiKey},
		"units": []string{"imperial"},
	}
	res, err := http.Get("https://pro.openweathermap.org/data/2.5/forecast?" + qs.Encode())
	if err != nil {
		return nil, err
	}

	var body owmResponse
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (f *openWeatherMap) NeedsRefresh() bool {
	return f.lastPull.Add(pullInterval).Before(time.Now())
}
