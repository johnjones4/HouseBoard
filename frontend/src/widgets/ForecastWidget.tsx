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

const days = [
  'Sun',
  'Mon',
  'Tue',
  'Wed',
  'Thu',
  'Fri',
  'Sat',
];

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
            speed: {
              type: 'linear' as const,
              display: true,
              position: 'left' as const,
              afterTickToLabelConversion: (axis: Scale) => {
                axis.ticks.forEach(tick => {
                  tick.label = `${tick.label}mph`
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
          labels: props.info.openWeatherMap.forecast.map((item, i) => {
            if (i === 0 || props.info.openWeatherMap.forecast[i - 1].timestamp.getDate() !== item.timestamp.getDate()) {
              return `${days[item.timestamp.getDay()]} ${item.timestamp.toLocaleTimeString()}`
            } else {
              return item.timestamp.toLocaleTimeString()
            }
          }),
          datasets: [
            {
              label: 'Temperature',
              data: props.info.openWeatherMap.forecast.map(item => {
                return item.temp
              }),
              borderColor: '#C05746',
              backgroundColor: '#C05746',
              yAxisID: 'temp',
              tension: 0.2,
            },
            {
              label: 'Precip Chance',
              data: props.info.openWeatherMap.forecast.map(item => {
                return item.pop
              }),
              borderColor: '#3E78B2',
              backgroundColor: '#3E78B2',
              yAxisID: 'pcnt',
              tension: 0.2,
            },
            {
              label: 'Humidity',
              data: props.info.openWeatherMap.forecast.map(item => {
                return item.humidity
              }),
              borderColor: '#0C7C59',
              backgroundColor: '#0C7C59',
              yAxisID: 'pcnt',
              tension: 0.2,
            },
            {
              label: 'Wind Speed',
              data: props.info.openWeatherMap.forecast.map(item => {
                return item.windSpeed
              }),
              borderColor: '#FEE440',
              backgroundColor: '#FEE440',
              yAxisID: 'speed',
              tension: 0.2,
            },
            {
              label: 'Cloud Coverage',
              data: props.info.openWeatherMap.forecast.map(item => {
                return item.cloudCov
              }),
              borderColor: '#800080',
              backgroundColor: '#800080',
              yAxisID: 'pcnt',
              tension: 0.2,
            },
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
