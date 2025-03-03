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

	if len(c.services.ICal.Events) > 0 {
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

	if len(c.services.Files.Files) > 0 {
		resp.Files = &Files{
			Files: make([]string, len(c.services.Files.Files)),
		}
		for i, file := range c.services.Files.Files {
			resp.Files.Files[i] = file
		}
	}

	if c.services.NOAA.Weather != nil {
		resp.Forecast = &Forecast{
			RadarURL: c.services.NOAA.Weather.RadarURL,
			Alerts:   make([]string, len(c.services.NOAA.Weather.Alerts)),
		}
		for i, alert := range c.services.NOAA.Weather.Alerts {
			resp.Forecast.Alerts[i] = alert.Headline
		}
	}

	if c.services.OpenWeatherMap.OwmResponse != nil {
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

	if c.services.WeatherStation != nil {
		resp.WeatherStation = &WeatherStation{
			AnemometerAverage: c.services.WeatherStation.WeatherStatonResponse.AvgWindSpeed,
			AnemometerMin:     c.services.WeatherStation.WeatherStatonResponse.MinWindSpeed,
			AnemometerMax:     c.services.WeatherStation.WeatherStatonResponse.MaxWindSpeed,
			Gas:               c.services.WeatherStation.WeatherStatonResponse.Gas,
			Pressure:          c.services.WeatherStation.WeatherStatonResponse.Pressure,
			RelativeHumidity:  c.services.WeatherStation.WeatherStatonResponse.RelativeHumidity,
			Temperature:       c.services.WeatherStation.WeatherStatonResponse.Temperature,
			Timestamp:         c.services.WeatherStation.WeatherStatonResponse.Timestamp,
			VaneDirection:     c.services.WeatherStation.WeatherStatonResponse.VaneDirection,
		}
	}

	if c.services.Traffic.TrafficResp != nil {
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

	if len(c.services.Trello.List) > 0 {
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

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		c.log.Error("error encoding response", slog.Any("error", err))
	}
}
