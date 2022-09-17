import React from 'react'
import { Info } from '../server/types'
import Widget from './Widget'

interface WeatherStationWidgetProps {
  info: Info
}

const WeatherStationWidget = (props: WeatherStationWidgetProps) => {
  return (
    <Widget name='weather-station' title='Weather Station'>
      
    </Widget>
  )
}

export default WeatherStationWidget
