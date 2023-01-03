import React, { useEffect, useState } from 'react'
import { Info } from '../server/types'
import Widget from './Widget'
import './ClockWidget.css'

interface ClockWidgetProps {
  info: Info
}

const months = [
  'jan',
  'feb',
  'mar',
  'apr',
  'may',
  'jun',
  'jul',
  'aug',
  'sep',
  'oct',
  'nov',
  'dec'
]

const prefixZeros = (n: number, z: number): string => {
  let ns = `${n}`
  while (ns.length < z) {
    ns = '0' + ns
  }
  return ns
}

const curStamp = (time: Date): string => {
  return `${months[time.getMonth()]} ${time.getDate()} ${prefixZeros(time.getHours(), 2)}:${prefixZeros(time.getMinutes(), 2)}:${prefixZeros(time.getSeconds(), 2)}`
}

const ClockWidget = (props: ClockWidgetProps) => {
  const [time, setTime] = useState(new Date())
  useEffect(() => {
    setInterval(() => setTime(new Date()), 500)
  },[])
  return (
    <Widget name='clock' title='Clock'>
      <div className='clock-str'>{curStamp(time)}</div>
    </Widget>
  )
}

export default ClockWidget
