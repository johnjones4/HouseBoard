export interface ICalEventResponse {
  title: string
  start: string
  end: string
  label: string
}

export type ICalResponse = ICalEventResponse[]

export interface NOAAForecastItemResponse {
  startTime: string
  endTime: string
  detailedForecast: string
  name: string
  temperature: number
  temperatureUnit: string
  windSpeed: string
  windDirection: string
  icon: string
  isDaytime: boolean
}

export interface NOAAResponse {
  radarURL: string
  forecast: NOAAForecastItemResponse[]
}

export interface WeatherStationResponse {
  timestamp: string
  anemometerMax: number
  temperature: number
  gas: number
  relativeHumidity: number
  pressure: number
  vaneDirection: number
}

export interface TrafficDestinationResponse {
  destination: string
  expectedDuration: number
  estimatedDuration: number
}

export interface TrafficResponse {
  destinations: TrafficDestinationResponse[]
}

export interface TrelloCardResponse {
  name: string
  id: string
}

export interface TrelloListResponse {
  cards: TrelloCardResponse[]
  name: string
}

export interface InfoResponse {
  ical: ICalResponse
  noaa: NOAAResponse
  weatherstation: WeatherStationResponse
  traffic: TrafficResponse
  trello: TrelloListResponse[]
}
