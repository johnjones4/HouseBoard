import React, { useEffect, useState } from 'react'
import { Info } from '../server/types'
import Widget from './Widget'
import './ClockWidget.css'

interface ClockWidgetProps {
  info: Info
}

const ClockWidget = (props: ClockWidgetProps) => {
  const [time, setTime] = useState(new Date())
  useEffect(() => {
    setInterval(() => setTime(new Date()), 500)
  },[])
  return (
    <Widget name='clock' title='Clock'>
      <div className='clock-str'>{time.toLocaleString()}</div>
    </Widget>
  )
}

export default ClockWidget
