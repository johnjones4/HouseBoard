// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Card defines model for Card.
type Card struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Event defines model for Event.
type Event struct {
	End   time.Time `json:"end"`
	Label string    `json:"label"`
	Start time.Time `json:"start"`
	Title string    `json:"title"`
}

// Events defines model for Events.
type Events struct {
	Events []Event `json:"events"`
}

// Files defines model for Files.
type Files struct {
	Files []string `json:"files"`
}

// Forecast defines model for Forecast.
type Forecast struct {
	Alerts           []string         `json:"alerts"`
	Forecast         []ForecastPeriod `json:"forecast"`
	LocalRadarURL    string           `json:"localRadarURL"`
	NationalRadarURL string           `json:"nationalRadarURL"`
}

// ForecastPeriod defines model for ForecastPeriod.
type ForecastPeriod struct {
	Datetime                   time.Time `json:"datetime"`
	FeelsLike                  float64   `json:"feelsLike"`
	Pressure                   float64   `json:"pressure"`
	ProbabilityOfPrecipitation float64   `json:"probabilityOfPrecipitation"`
	RelativeHumidity           float64   `json:"relativeHumidity"`
	Temperature                float64   `json:"temperature"`
	WindSpeed                  float64   `json:"windSpeed"`
}

// Info defines model for Info.
type Info struct {
	Events         *Events         `json:"events,omitempty"`
	Files          *Files          `json:"files,omitempty"`
	Forecast       *Forecast       `json:"forecast,omitempty"`
	Summary        *Summary        `json:"summary,omitempty"`
	SunriseSunset  *SunriseSunset  `json:"sunriseSunset,omitempty"`
	Traffic        *Traffic        `json:"traffic,omitempty"`
	Trello         *Trello         `json:"trello,omitempty"`
	Trivia         *Trivia         `json:"trivia,omitempty"`
	WeatherStation *WeatherStation `json:"weatherStation,omitempty"`
}

// List defines model for List.
type List struct {
	Cards []Card `json:"cards"`
	Name  string `json:"name"`
}

// Summary defines model for Summary.
type Summary struct {
	Summary string `json:"summary"`
}

// SunriseSunset defines model for SunriseSunset.
type SunriseSunset struct {
	Sunrise *time.Time `json:"sunrise,omitempty"`
	Sunset  *time.Time `json:"sunset,omitempty"`
}

// Traffic defines model for Traffic.
type Traffic struct {
	Destinations []TrafficDestination `json:"destinations"`
}

// TrafficDestination defines model for TrafficDestination.
type TrafficDestination struct {
	Destination string `json:"destination"`

	// EstimatedDuration Estimated travel duration in seconds based on current conditions.
	EstimatedDuration int `json:"estimatedDuration"`

	// ExpectedDuration Expected travel duration in seconds.
	ExpectedDuration int `json:"expectedDuration"`
}

// Trello defines model for Trello.
type Trello struct {
	List []List `json:"list"`
}

// Trivia defines model for Trivia.
type Trivia struct {
	Choices          []string `json:"choices"`
	PreviousAnswer   string   `json:"previousAnswer"`
	PreviousQuestion string   `json:"previousQuestion"`
	Question         string   `json:"question"`
}

// WeatherStation defines model for WeatherStation.
type WeatherStation struct {
	// AnemometerAverage Average wind speed measured by the anemometer.
	AnemometerAverage float64 `json:"anemometerAverage"`

	// AnemometerMax Maximum wind speed measured by the anemometer.
	AnemometerMax float64 `json:"anemometerMax"`

	// AnemometerMin Minimum wind speed measured by the anemometer.
	AnemometerMin float64 `json:"anemometerMin"`

	// Gas Gas concentration level measured.
	Gas float64 `json:"gas"`

	// Pressure Atmospheric pressure in hPa.
	Pressure float64 `json:"pressure"`

	// RelativeHumidity Relative humidity percentage.
	RelativeHumidity float64 `json:"relativeHumidity"`

	// Temperature Current temperature in degrees Celsius.
	Temperature float64 `json:"temperature"`

	// Timestamp The timestamp of the weather station data.
	Timestamp time.Time `json:"timestamp"`

	// VaneDirection Wind direction in degrees from the vane.
	VaneDirection float64 `json:"vaneDirection"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /info)
	GetInfo(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// (GET /info)
func (_ Unimplemented) GetInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetInfo operation middleware
func (siw *ServerInterfaceWrapper) GetInfo(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetInfo(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/info", wrapper.GetInfo)
	})

	return r
}
