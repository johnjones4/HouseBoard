package service

import "main/core"

type Configuration struct {
	ICal           ICalConfiguration           `json:"ical"`
	NOAA           NOAAConfiguration           `json:"noaa"`
	WeatherStation WeatherStationConfiguration `json:"weatherStation"`
}

func (c Configuration) Configurations() []core.ServiceConfig {
	return []core.ServiceConfig{
		c.ICal,
		c.NOAA,
		c.WeatherStation,
	}
}
