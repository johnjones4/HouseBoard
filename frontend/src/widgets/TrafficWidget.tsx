import React from 'react'
import { Info } from '../server/types'
import Widget from './Widget'

interface TrafficWidgetProps {
  info: Info
}

const TrafficWidget = (props: TrafficWidgetProps) => {
  return (
    <Widget name='traffic' title='Traffic'>
      
    </Widget>
  )
}

export default TrafficWidget
