import React from 'react'
import { Info, TrafficDestination } from '../server/types'
import { fancyTimeFormat } from '../util'
import Widget from './Widget'
import './TrafficWidget.css'

interface TrafficWidgetProps {
  info: Info
}

const timeClass = (d: TrafficDestination): string => {
  const pcnt = d.estimatedDuration / d.expectedDuration
  if (pcnt > 1.5) {
    return 'traffic-time-severe'
  } else if (pcnt > 1.1) {
    return 'traffic-time-delayed'
  }
  return 'traffic-time-normal'
}

const TrafficWidget = (props: TrafficWidgetProps) => {
  return (
    <Widget name='traffic' title='Traffic'>
      <table>
        <thead>
          <tr>
            <th>Desintation</th>
            <th>Estimate Time</th>
          </tr>
        </thead>
        <tbody>
          {props.info.traffic.destinations.map(d => (
            <tr key={d.destination}>
              <td>{d.destination}</td>
              <td className={'traffic-time ' + timeClass(d)}>{fancyTimeFormat(d.estimatedDuration)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </Widget>
  )
}

export default TrafficWidget
