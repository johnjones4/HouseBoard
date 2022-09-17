import React from 'react'
import { Info } from '../server/types'
import Widget from './Widget'
import './RadarWidget.css'

interface RadarWidgetProps {
  info: Info
}

const RadarWidget = (props: RadarWidgetProps) => {
  return (
    <Widget name='radar' title='Radar'>
      <div className='radar-image' style={{backgroundImage: `url(${props.info.noaa.radarURL})`}} />
    </Widget>
  )
}

export default RadarWidget
