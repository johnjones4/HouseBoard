import React from 'react'
import { Info } from '../server/types'
import Widget from './Widget'
import './WeatherStationWidget.css'

interface WeatherStationWidgetProps {
  info: Info
}

const directions = [
  "n",
  "ne",
  "e",
  "se",
  "s",
  "sw",
  "w",
  "nw"
]

const vaneDegreeToWindDirection = (ang: number): string|undefined => {
  const arc = 360.0 / directions.length
  return directions.find((_, i) => {
    const angle = arc * i
    let arcStart = angle - (arc/2)
    let arcEnd = angle + (arc/2) % 360
    if (arcStart < 0) {
      arcStart = 360 - angle - (arc/2)
    }
    return ang >= arcStart && ang <= arcEnd
  })
}

const WeatherStationWidget = (props: WeatherStationWidgetProps) => {
  if (!props.info.weatherstation) {
    return null
  }
  const direction = vaneDegreeToWindDirection(props.info.weatherstation.vaneDirection)
  const tableInfo = [
    {
      name: 'Wind',
      value: `${props.info.weatherstation.anemometerMax.toFixed(2)} mph`,
      classes: [`weather-station-wind-${direction}`]
    },
    {
      name: 'Humidity',
      value: `${props.info.weatherstation.relativeHumidity.toFixed(0)} %`,
      classes: []
    },
    {
      name: 'Pressure',
      value: `${props.info.weatherstation.pressure.toFixed(1)} inHg`,
      classes: []
    }
  ]
  return (
    <Widget name='weather-station' title='Weather Station'>
      <div className='weather-station-temperature'>
        {props.info.weatherstation.temperature.toFixed(1)}&deg;
      </div>
      <div className='weather-station-details'>
        <table>
          <tbody>
            {tableInfo.map((row) => (
              <tr key={row.name}>
                <td>{row.name}</td>
                <td className='weather-station-detail-value'>
                  <span className={row.classes.join(' ')}>
                    {row.value}  
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Widget>
  )
}

export default WeatherStationWidget
