package api

import (
	"encoding/json"
	"log/slog"
	"main/service"
	"net/http"
	"time"
)

type controller struct {
	services *service.Services
	log      *slog.Logger
}

func (c *controller) GetInfo(w http.ResponseWriter, r *http.Request) {
	c.services.Lock.RLock()
	defer c.services.Lock.RUnlock()

	resp := Info{
		Summary: &Summary{
			Summary: c.services.Summary(),
		},
	}

	if c.services.ICal != nil && len(c.services.ICal.Events) > 0 {
		resp.Events = &Events{
			Events: make([]Event, len(c.services.ICal.Events)),
		}
		for i, event := range c.services.ICal.Events {
			resp.Events.Events[i] = Event{
				End:   event.End,
				Start: event.Start,
				Label: event.Label,
				Title: event.Title,
			}
		}
	}

	if c.services.Files != nil && len(c.services.Files.Files) > 0 {
		resp.Files = &Files{
			Files: c.services.Files.Files,
		}
	}

	if c.services.NOAA != nil && c.services.NOAA.Weather != nil {
		resp.Forecast = &Forecast{
			LocalRadarURL:    c.services.NOAA.Weather.LocalRadarURL,
			NationalRadarURL: c.services.NOAA.Weather.NationalRadarURL,
			Alerts:           make([]string, len(c.services.NOAA.Weather.Alerts)),
		}
		for i, alert := range c.services.NOAA.Weather.Alerts {
			resp.Forecast.Alerts[i] = alert.Headline
		}
	}

	if c.services.OpenWeatherMap != nil && c.services.OpenWeatherMap.OwmResponse != nil {
		if resp.Forecast == nil {
			resp.Forecast = &Forecast{}
		}
		resp.Forecast.Forecast = make([]ForecastPeriod, len(c.services.OpenWeatherMap.OwmResponse.List))
		for i, period := range c.services.OpenWeatherMap.OwmResponse.List {
			resp.Forecast.Forecast[i] = ForecastPeriod{
				Datetime:                   time.Unix(period.Dt, 0),
				FeelsLike:                  period.Main.FeelsLike,
				ProbabilityOfPrecipitation: period.Pop,
				RelativeHumidity:           float64(period.Main.Humidity) / 100,
				Temperature:                period.Main.Temp,
				WindSpeed:                  period.Wind.Speed,
				Pressure:                   float64(period.Main.Pressure),
			}
		}
	}

	if c.services.WeatherStation != nil && c.services.WeatherStation.WeatherStatonResponse != nil {
		resp.WeatherStation = &WeatherStation{
			WindSpeed:     *c.services.WeatherStation.WeatherStatonResponse.WindSpeed,
			Gas:           *c.services.WeatherStation.WeatherStatonResponse.Gas,
			Pressure:      *c.services.WeatherStation.WeatherStatonResponse.Pressure,
			Humidity:      *c.services.WeatherStation.WeatherStatonResponse.Humidity,
			Temperature:   *c.services.WeatherStation.WeatherStatonResponse.Temperature,
			Timestamp:     *c.services.WeatherStation.LastUpdated,
			VaneDirection: *c.services.WeatherStation.WeatherStatonResponse.VaneDirection,
			Rainfall:      *c.services.WeatherStation.WeatherStatonResponse.Rainfall,
		}
	}

	if c.services.Traffic != nil && c.services.Traffic.TrafficResp != nil {
		resp.Traffic = &Traffic{
			Destinations: make([]TrafficDestination, len(c.services.Traffic.TrafficResp.Destinations)),
		}
		for i, dest := range c.services.Traffic.TrafficResp.Destinations {
			resp.Traffic.Destinations[i] = TrafficDestination{
				Destination:       dest.Destination,
				EstimatedDuration: dest.EstimatedDuration,
				ExpectedDuration:  dest.ExpectedDuration,
			}
		}
	}

	if c.services.Trello != nil && len(c.services.Trello.List) > 0 {
		resp.Trello = &Trello{
			List: make([]List, len(c.services.Trello.List)),
		}
		for i, list := range c.services.Trello.List {
			resp.Trello.List[i] = List{
				Name:  list.Name,
				Cards: make([]Card, len(list.Cards)),
			}
			for j, card := range list.Cards {
				resp.Trello.List[i].Cards[j] = Card{
					Id:   card.Id,
					Name: card.Name,
				}
			}
		}
	}

	if c.services.SunriseSunset != nil {
		resp.SunriseSunset = &SunriseSunset{
			Sunrise: c.services.SunriseSunset.Rise,
			Sunset:  c.services.SunriseSunset.Set,
		}
	}

	if c.services.Trivia != nil && c.services.Trivia.Current != nil && c.services.Trivia.Previous != nil {
		resp.Trivia = &Trivia{
			Question: c.services.Trivia.Current.Question,
			Choices:  c.services.Trivia.Current.Choices,
		}
		if c.services.Trivia.Previous.Answer >= 0 && c.services.Trivia.Previous.Answer < len(c.services.Trivia.Previous.Choices) {
			resp.Trivia.PreviousQuestion = c.services.Trivia.Previous.Question
			resp.Trivia.PreviousAnswer = c.services.Trivia.Previous.Choices[c.services.Trivia.Previous.Answer]
		}
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		c.log.Error("error encoding response", slog.Any("error", err))
	}
}
