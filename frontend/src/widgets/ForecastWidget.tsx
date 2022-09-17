import React from 'react'
import { Info } from '../server/types'
import Widget from './Widget'
import './ForecastWidget.css'

interface ForecastWidgetProps {
  info: Info
}

const ForecastWidget = (props: ForecastWidgetProps) => {
  return (
    <Widget name='forecast' title='Forecast'>
      { props.info.noaa.forecast.slice(0, 5).map((f, i) => {
        return (
          <div key={i} className='forecast-item'>
            <div className='forecast-title'>{f.name }</div>
            <div className='forecast-details'>
              <img className='forecast-icon' src={f.icon} alt={f.detailedForecast} />
              <div className='forecast-stats'>
                <div className='forecast-temp'>{f.temperature}°</div>
                <div className={`forecast-wind forecast-wind-${f.windDirection.toLowerCase()}`}>{f.windSpeed}</div>
              </div>
            </div>
            <div className='forecast-description'>{f.detailedForecast}</div>
          </div>
        )
      }) }
    </Widget>
  )
}

export default ForecastWidget
