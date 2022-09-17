import React from 'react'

export interface WidgetProps {
  title: string
  children?: any
  name: string
}

const Widget = (props: WidgetProps) => {
  return (
    <div className={`widget widget-${props.name}`}>
      <div className="widget-title">{props.title}</div>
      <div className="widget-body">
        {props.children}
      </div>
    </div>
  )
}

export default Widget
