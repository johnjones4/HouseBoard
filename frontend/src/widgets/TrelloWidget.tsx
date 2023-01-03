import React from 'react'
import { Info } from '../server/types'
import { fancyTimeFormat } from '../util'
import Widget from './Widget'
import './TrafficWidget.css'

interface TrelloWidgetProps {
  info: Info
}

const TrelloWidget = (props: TrelloWidgetProps) => {
  return (
    <Widget name='trello' title='Trello'>
      { props.info.trello.map(t => (
        <div className='trello-list' key={t.name}>
          <table>
            <tbody>
              {t.cards.slice(0,7).map(c => (
                <tr key={c.id}>
                  <td>
                    {c.name}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )) }
    </Widget>
  )
}

export default TrelloWidget
