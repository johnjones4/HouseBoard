import React, { useEffect, useState } from 'react';
import { getInfo } from './server/comm';
import { Info } from './server/types';
import Widget from './widgets/Widget';
import './App.css';
import CalendarWidget from './widgets/CalendarWidget';
import TrafficWidget from './widgets/TrafficWidget';
import ForecastWidget from './widgets/ForecastWidget';
import RadarWidget from './widgets/RadarWidget';
import WeatherStationWidget from './widgets/WeatherStationWidget';
import ClockWidget from './widgets/ClockWidget';
import TrelloWidget from './widgets/TrelloWidget';

function App() {
  const [info, setInfo] = useState(undefined as Info | undefined)

  useEffect(() => {
    const callServer = () => {
      getInfo()
      .then(i => {
        setInfo(i)
      })
      .catch(e => console.error(e))
    }
    callServer()
    setInterval(callServer, 1000 * 60 * 5)
  }, [])

  if (!info) {
    return null
  }
  return (
    <div className="dashboard">
      <CalendarWidget info={info} />
      <TrafficWidget info={info} />
      <ForecastWidget info={info} />
      <RadarWidget info={info} />
      <WeatherStationWidget info={info} />
      <ClockWidget info={info} />
      <TrelloWidget info={info} />
    </div>
  )
}

export default App;
