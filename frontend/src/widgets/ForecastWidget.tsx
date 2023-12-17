import React from 'react'
import { Info } from '../server/types'
import Widget from './Widget'
import './ForecastWidget.css'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Scale,
} from 'chart.js';
import { Line } from 'react-chartjs-2';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

// class MyScale extends Scale {
//   /* extensions ... */
//   id = 'myscale'
//   defaults = 
// }
// MyScale.id = 'myScale';
// MyScale.defaults = defaultConfigObject;

interface ForecastWidgetProps {
  info: Info
}

const ForecastWidget = (props: ForecastWidgetProps) => {
  return (
    <Widget name='forecast' title='Forecast'>
      <Line 
        options={{
          aspectRatio: 6,
          responsive: true,
          maintainAspectRatio: true,
          plugins: {
            legend: {
              display: true
            },
            title: {
              display: false,
            },
          },
          scales: {
            temp: {
              type: 'linear' as const,
              display: true,
              position: 'left' as const,
              afterTickToLabelConversion: (axis: Scale) => {
                axis.ticks.forEach(tick => {
                  tick.label = `${tick.label}°`
                })
              }
            },
            pcnt: {
              type: 'linear' as const,
              display: true,
              position: 'left' as const,
              min: 0,
              max: 100,
              grid: {
                display: true,
                color: '#202020'
              },
              afterTickToLabelConversion: (axis: Scale) => {
                axis.ticks.forEach(tick => {
                  tick.label = `${tick.label}%`
                })
              }
            },
          },
        }} 
        data={{
          labels: props.info.noaa.forecast.map(item => {
            return item.name
          }),
          datasets: [
            {
              label: 'Temperature',
              data: props.info.noaa.forecast.map(item => {
                return item.temperature
              }),
              borderColor: '#C05746',
              backgroundColor: '#C05746',
              yAxisID: 'temp',
              tension: 0.2,
            },
            {
              label: 'Precip Chance',
              data: props.info.noaa.forecast.map(item => {
                return item.probabilityOfPrecipitation.value
              }),
              borderColor: '#3E78B2',
              backgroundColor: '#3E78B2',
              yAxisID: 'pcnt',
              tension: 0.2,
            },
            {
              label: 'Humidity',
              data: props.info.noaa.forecast.map(item => {
                return item.relativeHumidity.value
              }),
              borderColor: '#0C7C59',
              backgroundColor: '#0C7C59',
              yAxisID: 'pcnt',
              tension: 0.2,
            },
            {
              label: 'Dewpoint',
              data: props.info.noaa.forecast.map(item => {
                return (item.dewpoint.value * 1.8) + 32
              }),
              borderColor: '#FEE440',
              backgroundColor: '#FEE440',
              yAxisID: 'temp',
              tension: 0.2,
            }
          ]
        }
      } />
      {/* { props.info.noaa.forecast.slice(0, 5).map((f, i) => {
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
      }) } */}
    </Widget>
  )
}

export default ForecastWidget
