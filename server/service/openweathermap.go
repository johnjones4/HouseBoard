package service

import (
	"context"
	"encoding/json"
	"fmt"
	"main/core"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type OwmResponse struct {
	List []OwmResponseItem `json:"list"`
}

type OwmResponseItem struct {
	Dt         int64        `json:"dt"`
	Main       OwmMain      `json:"main"`
	Weather    []OwmWeather `json:"weather"`
	Clouds     OwmClouds    `json:"clouds"`
	Wind       OwmWind      `json:"wind"`
	Visibility int          `json:"visibility"`
	Pop        float64      `json:"pop"`
	Sys        OwmSys       `json:"sys"`
	DtTxt      string       `json:"dt_txt"`
}

// Main represents the main temperature details.
type OwmMain struct {
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
type OwmWeather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// Clouds represents the cloud coverage details.
type OwmClouds struct {
	All int `json:"all"`
}

// Wind represents the wind details.
type OwmWind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

// Sys represents the system details.
type OwmSys struct {
	Pod string `json:"pod"`
}

type owmConfiguration struct {
	Location core.Coordinate
	ApiKey   string `json:"apiKey"`
}

func (c owmConfiguration) Empty() bool {
	return c.Location.Latitude == 0 && c.Location.Longitude == 0
}

func (c owmConfiguration) Service() *OpenWeatherMap {
	return &OpenWeatherMap{
		owmConfiguration: c,
	}
}

const pullInterval = time.Hour

type OpenWeatherMap struct {
	owmConfiguration
	lastPull    time.Time
	OwmResponse *OwmResponse
}

func (f *OpenWeatherMap) Name() string {
	return "openWeatherMap"
}

func (f *OpenWeatherMap) Refresh(c context.Context) error {
	qs := url.Values{
		"lat":   []string{fmt.Sprint(f.Location.Latitude)},
		"lon":   []string{fmt.Sprint(f.Location.Longitude)},
		"appid": []string{f.ApiKey},
		"units": []string{"imperial"},
	}
	url := "https://pro.openweathermap.org/data/2.5/forecast?" + qs.Encode()
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	var body OwmResponse
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return err
	}

	f.OwmResponse = &body

	return nil
}

func (f *OpenWeatherMap) NeedsRefresh() bool {
	if f == nil {
		return false
	}
	return f.lastPull.Add(pullInterval).Before(time.Now())
}

func (f *OpenWeatherMap) StateForPrompt() *string {
	if f == nil {
		return nil
	}

	if f.OwmResponse == nil {
		return nil
	}

	var summary strings.Builder

	summary.WriteString("Weather Forecast:\n")

	for _, item := range f.OwmResponse.List {
		summary.WriteString(fmt.Sprintf(`%s:
- Temperature: %0.2f degrees Fahrenheit
- Feels Like Temp: %0.2f degrees Fahrenheit
- Wind: %0.2f MPH
- Change of Precipitation: %0.2f%%
- Pressure: %d inHg
- Humidity: %d%%`,
			time.Unix(item.Dt, 0).String(),
			item.Main.Temp,
			item.Main.FeelsLike,
			item.Wind.Speed,
			item.Pop*100,
			item.Main.Pressure,
			item.Main.Humidity,
		))
		summary.WriteString("\n")
	}

	str := summary.String()

	return &str
}
