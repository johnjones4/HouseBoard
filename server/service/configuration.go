package service

type Configuration struct {
	ICal           iCalConfiguration           `json:"ical"`
	NOAA           noaaConfiguration           `json:"noaa"`
	WeatherStation weatherStationConfiguration `json:"weatherStation"`
	Traffic        trafficConfiguration        `json:"traffic"`
	Trello         trelloConfiguration         `json:"trello"`
	Files          FileConfiguration           `json:"file"`
	OpenWeatherMap OpenWeatherMap              `json:"openWeatherMap"`
	SunriseSunset  SunriseSunsetConfiguration  `json:"sunriseSunset"`
	Trivia         TriviaConfiguration         `json:"trivia"`
}

func (c Configuration) Services() *Services {
	svc := &Services{}
	if !c.ICal.Empty() {
		svc.ICal = c.ICal.Service()
	}
	if !c.NOAA.Empty() {
		svc.NOAA = c.NOAA.Service()
	}
	if !c.OpenWeatherMap.Empty() {
		svc.OpenWeatherMap = c.OpenWeatherMap.Service()
	}
	if !c.Traffic.Empty() {
		svc.Traffic = c.Traffic.Service()
	}
	if !c.Trello.Empty() {
		svc.Trello = c.Trello.Service()
	}
	if !c.Files.Empty() {
		svc.Files = c.Files.Service()
	}
	if !c.WeatherStation.Empty() {
		svc.WeatherStation = c.WeatherStation.Service()
	}
	if !c.SunriseSunset.Empty() {
		svc.SunriseSunset = c.SunriseSunset.Service()
	}
	if !c.Trivia.Empty() {
		svc.Trivia = c.Trivia.Service()
	}

	return svc
}
