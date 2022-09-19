package service

import "main/core"

type Configuration struct {
	ICal           iCalConfiguration           `json:"ical"`
	NOAA           noaaConfiguration           `json:"noaa"`
	WeatherStation weatherStationConfiguration `json:"weatherStation"`
	Traffic        trafficConfiguration        `json:"traffic"`
}

func (c Configuration) Configurations() []core.ServiceConfig {
	return []core.ServiceConfig{
		c.ICal,
		c.NOAA,
		c.WeatherStation,
		c.Traffic,
	}
}
