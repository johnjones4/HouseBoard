import { ICalEventResponse, ICalResponse, InfoResponse, NOAAForecastItemResponse, NOAAResponse, TrafficDestinationResponse, TrafficResponse, WeatherStationResponse } from "./responses"

export class ICalEvent {
  title: string
  start: Date
  end: Date
  label: string

  constructor(r: ICalEventResponse) {
    this.title = r.title
    this.start = new Date(Date.parse(r.start))
    this.end = new Date(Date.parse(r.end))
    this.label = r.label
  }
}

export class ICal {
  calendars: ICalEvent[]
  labels: string[]

  constructor(r: ICalResponse) {
    this.labels = []
    this.calendars = r.map(r => {
      if (this.labels.indexOf(r.label) < 0) {
        this.labels.push(r.label)
      }
      return new ICalEvent(r)
    })
  }
}

export class NOAAForecastItem {
  startTime: Date
  endTime: Date
  detailedForecast: string
  name: string
  temperature: number
  temperatureUnit: string
  windSpeed: string
  windDirection: string 
  icon: string
  isDaytime: boolean

  constructor(r: NOAAForecastItemResponse) {
    this.startTime = new Date(Date.parse(r.startTime))
    this.endTime = new Date(Date.parse(r.endTime))
    this.detailedForecast = r.detailedForecast
    this.name = r.name
    this.temperature = r.temperature
    this.temperatureUnit = r.temperatureUnit
    this.windSpeed = r.windSpeed
    this.windDirection = r.windDirection
    this.icon = r.icon
    this.isDaytime = r.isDaytime
  }
}


export class NOAA {
  radarURL: string
  forecast: NOAAForecastItem[]

  constructor(r: NOAAResponse) {
    this.radarURL = r.radarURL
    this.forecast = r.forecast.map(i => new NOAAForecastItem(i))
  }
}

export class WeatherStation {
  timestamp: Date
  uptime: number
  avgWindSpeed: number
  minWindSpeed: number
  maxWindSpeed: number
  temperature: number
  gas: number
  relativeHumidity: number
  pressure: number

  constructor(r: WeatherStationResponse) {
    this.timestamp = new Date(Date.parse(r.timestamp))
    this.uptime = r.uptime
    this.avgWindSpeed = r.avg_wind_speed
    this.minWindSpeed = r.min_wind_speed
    this.maxWindSpeed = r.max_wind_speed
    this.temperature = r.temperature
    this.gas = r.gas
    this.relativeHumidity = r.relative_humidity
    this.pressure = r.pressure
  }
}

export class TrafficDestination {
  destination: string
  expectedDuration: number
  estimatedDuration: number

  constructor(r: TrafficDestinationResponse) {
    this.destination = r.destination
    this.expectedDuration = r.expectedDuration
    this.estimatedDuration = r.estimatedDuration
  }
}

export class Traffic {
  destinations: TrafficDestination[]

  constructor(r: TrafficResponse) {
    this.destinations = r.destinations.map(d => new TrafficDestination(d))
  }
}

export class Info {
  ical: ICal
  noaa: NOAA
  weatherstation: WeatherStation
  traffic: Traffic
  
  constructor(r: InfoResponse) {
    this.ical = new ICal(r.ical)
    this.noaa = new NOAA(r.noaa)
    this.weatherstation = new WeatherStation(r.weatherstation)
    this.traffic = new Traffic(r.traffic)
  }
}
