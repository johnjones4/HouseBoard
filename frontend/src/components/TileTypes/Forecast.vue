<template>
  <Tile :position="props.position" name="forecast">
    <Line
      :options="chartOptions"
      :data="chartData"
    />
  </Tile>
</template>

<script lang="ts" setup>

import { useInfoStore } from '../../stores/info';
import type { Position } from '../Position';
import Tile from '../Tile.vue';
import { computed } from 'vue';
import { Line } from 'vue-chartjs';
import { Chart as ChartJS, Scale, Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement } from 'chart.js'
import type { ChartData, ChartOptions } from 'chart.js';
import { Color } from '../colors';

ChartJS.register(Title, Tooltip, Legend, LineElement, PointElement, CategoryScale, LinearScale);

const infoStore = useInfoStore();

const chartOptions : ChartOptions<'line'> = {
  responsive: true,
  maintainAspectRatio: false,
  animation: false,
  scales: {
    temp: {
      type: 'linear' as const,
      display: true,
      position: 'left' as const,
      afterTickToLabelConversion: (axis: Scale) => {
        axis.ticks.forEach(tick => {
          tick.label = `${tick.label}°`
        })
      }
    },
    pcnt: {
      type: 'linear' as const,
      display: true,
      position: 'left' as const,
      min: 0,
      max: 100,
      grid: {
        display: true,
        color: '#202020'
      },
      afterTickToLabelConversion: (axis: Scale) => {
        axis.ticks.forEach(tick => {
          tick.label = `${tick.label}%`
        })
      }
    },
  },
};

const props = defineProps<{
  position: Position;
}>();

const days = [
  'Sun',
  'Mon',
  'Tue',
  'Wed',
  'Thu',
  'Fri',
  'Sat'
]

const chartData = computed(():ChartData<'line', number[], unknown> => {
  if (!infoStore.info.forecast) {
    return {
      labels: [],
      datasets: [],
    };
  }
  const out = {
    labels: infoStore.info.forecast.forecast.map((f): string => {
      const date = new Date(f.datetime);
      const weekday = days[date.getDay()];
      return `${weekday} @ ${date.toLocaleTimeString()}`
    }),
    datasets: [
      {
        label: 'Temperature',
        data: infoStore.info.forecast.forecast.map(f => f.temperature),
        yAxisID: 'temp',
        borderColor: Color.orange,
        backgroundColor: Color.orange,
      },
      {
        label: 'Feels Like',
        data: infoStore.info.forecast.forecast.map(f => f.feelsLike),
        yAxisID: 'temp',
        borderColor: Color.yellow,
        backgroundColor: Color.yellow,
      },
      {
        label: 'Humidity',
        data: infoStore.info.forecast.forecast.map(f => f.relativeHumidity * 100.0),
        yAxisID: 'pcnt',
        borderColor: Color.green,
        backgroundColor: Color.green,
      },
      {
        label: 'Precipition Chance',
        data: infoStore.info.forecast.forecast.map(f => f.probabilityOfPrecipitation * 100.0),
        yAxisID: 'pcnt',
        borderColor: Color.blue,
        backgroundColor: Color.blue,
      },
    ]
  }
  return out;
});

</script>

<style>


</style>