export interface ICalEventResponse {
  title: string
  start: string
  end: string
  label: string
}

export type ICalResponse = ICalEventResponse[]

export interface NOAAForecastValue {
  value: number
  unitCode: string
}

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
  probabilityOfPrecipitation: NOAAForecastValue
  relativeHumidity: NOAAForecastValue
  dewpoint: NOAAForecastValue
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

export interface FileResponse {
  files: string[]
}

export interface OWMWeatherResponseBody {
  list: OWMWeatherResponse[]
}

// Main weather response structure
export interface OWMWeatherResponse {
  dt: number;
  main: OWMMain;
  weather: OWMWeather[];
  clouds: OWMClouds;
  wind: OWMWind;
  visibility: number;
  pop: number;
  sys: OWMSys;
  dt_txt: string;
}

// Main temperature details
export interface OWMMain {
  temp: number;
  feels_like: number;
  temp_min: number;
  temp_max: number;
  pressure: number;
  sea_level: number;
  grnd_level: number;
  humidity: number;
  temp_kf: number;
}

// Weather condition details
export interface OWMWeather {
  id: number;
  main: string;
  description: string;
  icon: string;
}

// Cloud coverage details
export interface OWMClouds {
  all: number;
}

// Wind details
export interface OWMWind {
  speed: number;
  deg: number;
  gust: number;
}

// System details
export interface OWMSys {
  pod: string;
}


export interface InfoResponse {
  ical: ICalResponse
  noaa: NOAAResponse
  weatherstation: WeatherStationResponse
  traffic: TrafficResponse
  trello: TrelloListResponse[]
  file: FileResponse
  openWeatherMap: OWMWeatherResponseBody
}
