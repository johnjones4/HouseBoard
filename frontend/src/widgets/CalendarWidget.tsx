import React from 'react'
import { ICal, ICalEvent, Info } from '../server/types'
import Widget from './Widget'
import './CalendarWidget.css'
import { hoursMinutesString } from '../util'

interface CalendarWidgetProps {
  info: Info
}

const daysOfWeek = ['Su', 'M', 'Tu', 'W', 'Th', 'F', 'Sa']

const timeString = (event: ICalEvent): string => {
  const start = hoursMinutesString(event.start)
  const end = hoursMinutesString(event.end)
  return start === end ? '' : `${start} | `
}

const CalendarWidget = (props: CalendarWidgetProps) => {
  const calendar = makeCalendarArray(props.info.ical.calendars)

  return (
    <Widget name='calendar' title='Calendar'>
      <div className='calendar-calendar'>
        {daysOfWeek.map(d => (<div className='calendar-header' key={d}>{d}</div>))}
        {calendar.map((c, i) => {
          if (!c.date) {
            return (<div className='calendar-item calendar-item-empty' key={i}/>)
          }
          return (
            <div key={i} className='calendar-item'>
              <div className='calendar-item-date'>
                {c.date.getDate()}
              </div>
              { c.events.length > 0 && (
                <ul className='calendar-item-events'>
                  {c.events.map((e, j) => {
                    const labelIndex = props.info.ical.labels.indexOf(e.label)
                    return (
                      <li className={`calendar-item-event calendar-item-event-${labelIndex}`} key={j}>
                        {timeString(e)}{e.title}
                      </li>
                    )
                  })}
                </ul>
              ) }
            </div>
          )
        })}
      </div>
    </Widget>
  )
}

interface CalendarItem {
  events: ICalEvent[]
  date: Date | null
}

const makeCalendarArray = (events: ICalEvent[]): CalendarItem[] => {
  const now = new Date()
  let curDay = now
  const array = [] as CalendarItem[]
  for (let row = 0; row < 5; row++) {
    if (curDay.getDay() > 0) {
      for (let i = 0; i < curDay.getDay(); i++) {
        array.push({
          events: [],
          date: null,
        })
      }
    }
    for (let col = curDay.getDay(); col < 7; col++) {      
      const date = curDay
      array.push({
        events: events.filter(event => eventOccursOnDay(event, date)),
        date: curDay
      })
      curDay = new Date(curDay.getTime() + (1000*60*60*24))
    }
  }
  return array
} 

const eventOccursOnDay = (event: ICalEvent, day: Date): boolean => {
  const startOfDay = new Date(day.getFullYear(), day.getMonth(), day.getDate(), 0, 0, 0, 0)
  const endOfDay = new Date(startOfDay.getTime() + (1000*60*60*24))
  return (event.start.getTime() >= startOfDay.getTime() && event.start.getTime() < endOfDay.getTime())
    || (event.end.getTime() > startOfDay.getTime() && event.end.getTime() <= endOfDay.getTime())
    || (startOfDay.getTime() >= event.start.getTime() && endOfDay.getTime() <= event.end.getTime())
}

export default CalendarWidget
