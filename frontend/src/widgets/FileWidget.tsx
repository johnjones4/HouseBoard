import React from 'react'
import { Info } from '../server/types'
import Widget from './Widget'
import './FileWidget.css'

interface FileWidgetProps {
  info: Info
}

const FileWidget = (props: FileWidgetProps) => {
  return (
    <Widget name='file' title='File'>
      <div className='file-image' style={{backgroundImage: `url(${props.info.file.files[0]})`}} />
    </Widget>
  )
}

export default FileWidget
