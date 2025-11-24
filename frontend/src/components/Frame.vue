<template>
  <div class="frame-outer">
    <div class="frame-title">{{ props.title }}</div>
    <div class="frame-inner">
      <template v-for="tile in props.tiles" :key="tile.tileType">
        <Files 
          v-if="tile.tileType === TileType.files"
          :position="tile.position"
        />
        <Calendar 
          v-else-if="tile.tileType === TileType.calendar" 
          :position="tile.position"
        />
        <Forecast 
          v-else-if="tile.tileType === TileType.forecast"
          :position="tile.position"
        />
        <Radar 
          v-else-if="tile.tileType === TileType.localRadar" 
          :position="tile.position"
          :radar-type="RadarType.local"
        />
        <Radar 
          v-else-if="tile.tileType === TileType.nationalRadar" 
          :position="tile.position"
          :radar-type="RadarType.national"
        />
        <Trello 
          v-else-if="tile.tileType === TileType.trello" 
          :position="tile.position"
        />
        <WeatherStation 
          v-else-if="tile.tileType === TileType.weatherStation" 
          :position="tile.position"
        />
        <Traffic 
          v-else-if="tile.tileType === TileType.traffic" 
          :position="tile.position"
        />
        <Clock 
          v-else-if="tile.tileType === TileType.clock" 
          :position="tile.position"
        />
        <Summary 
          v-else-if="tile.tileType === TileType.summary" 
          :position="tile.position"
        />
        <Agenda
          v-else-if="tile.tileType === TileType.agenda"
          :position="tile.position"
        />
        <Trivia
          v-else-if="tile.tileType === TileType.trivia"
          :position="tile.position"
        />
      </template>
    </div>
    <div v-if="infoStore && infoStore.info && infoStore.info.weatherStation" class="frame-footer">
      <div v-for="item in footers" :key="item.label" class="frame-footer-prop">
        <FontAwesomeIcon :icon="item.icon" size="3x" />
        <div class="frame-footer-prop-info">
          <div class="frame-footer-prop-value" :ref="(v) => setFooterItemRef(item, v as HTMLDivElement)">
            {{ item.value }}
          </div>
          <div class="frame-footer-prop-label">
            {{ item.label }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { TileType } from './TileType';
import Calendar from './TileTypes/Calendar.vue';
import Files from './TileTypes/Files.vue';
import Clock from './TileTypes/Clock.vue';
import Forecast from './TileTypes/Forecast.vue';
import Radar from './TileTypes/Radar.vue';
import { RadarType } from './RadarType';
import Summary from './TileTypes/Summary.vue';
import Traffic from './TileTypes/Traffic.vue';
import Trello from './TileTypes/Trello.vue';
import WeatherStation from './TileTypes/WeatherStation.vue';
import type { FrameProps } from './FrameProps';
import Agenda from './TileTypes/Agenda.vue';
import { useInfoStore } from '../stores/info';
import { computed, ref } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { faClock, faSun, faMoon } from '@fortawesome/free-regular-svg-icons'
import { faTemperatureHalf } from '@fortawesome/free-solid-svg-icons'
import { library, type IconDefinition } from '@fortawesome/fontawesome-svg-core'
import Trivia from './TileTypes/Trivia.vue';

library.add(faClock, faTemperatureHalf, faSun, faMoon);

const props = defineProps<FrameProps>();

const dateTimeEl = ref<null|HTMLDivElement>(null);

const infoStore = useInfoStore();

interface FooterSummaryItem {
  icon: IconDefinition;
  label: string;
  value: string;
}

const footers = computed((): FooterSummaryItem[] => {
  if (!infoStore.info) {
    return [];
  }
  return [
    {
      icon: faClock,
      label: 'Date/Time',
      value: new Date().toLocaleString(),
    },
    {
      icon: faTemperatureHalf,
      label: 'Temperature',
      value: infoStore.info.weatherStation ? `${infoStore.info.weatherStation.temperature.toFixed(1)}°` : '',
    },
    {
      icon: faSun,
      label: 'Sunrise',
      value: infoStore.info.sunriseSunset && infoStore.info.sunriseSunset.sunrise ? new Date(infoStore.info.sunriseSunset.sunrise).toLocaleTimeString() : '',
    },
    {
      icon: faMoon,
      label: 'Sunset',
      value: infoStore.info.sunriseSunset && infoStore.info.sunriseSunset.sunset ? new Date(infoStore.info.sunriseSunset.sunset).toLocaleTimeString() : '',
    },
  ]
});

const setFooterItemRef = (item: FooterSummaryItem, r: HTMLDivElement) => {
  if (item.label === 'Date/Time') {
    dateTimeEl.value = r;
  }
}

const tick = () => {
  if (dateTimeEl.value) {
    dateTimeEl.value.textContent = new Date().toLocaleString();
  }
  requestAnimationFrame(tick);
}

requestAnimationFrame(tick);

</script>

<style scoped>

.frame-outer {
  position: fixed;
  top: var(--default-padding);
  right: var(--default-padding);
  bottom: calc(var(--default-padding) + var(--progress-line-thickness));
  left: var(--default-padding);
  display: flex;
  flex-direction: column;
}

.frame-title {
  font-size: 1.25em;
  padding-bottom: 10px;
}

.frame-inner {
  display: grid;
  grid-template-columns: repeat(12, minmax(0, 1fr));
  grid-template-rows: repeat(6, minmax(0, 1fr));
  flex-grow: 1;
  grid-gap: var(--default-padding);
  padding: var(--default-padding);
  border: solid 1px var(--color-text-light);
  overflow: hidden;
}

.frame-footer {
  margin-top: var(--thin-padding);
  background-color: var(--color-text);
  color: var(--color-background);
  display: flex;
  flex-direction: row;
  padding: var(--thin-padding);
  font-size: 1.5;
}

.frame-footer-prop {
  margin-left: var(--default-padding);
  padding-left: var(--default-padding);
  border-left: dotted 1px var(--color-background);
  display: flex;
  flex-direction: row;
  align-items: center;
}

.frame-footer-prop-info {
  margin-left: var(--thin-padding);
}

.frame-footer-prop:first-child {
  margin-left: var(--thin-padding);
  border: none;
  padding: 0;
}

.frame-footer-prop-value {
  font-size: 2.5em;
  font-optical-sizing: auto;
  font-family: "Doto", monospace;
  font-weight: bold;
}

.frame-footer-prop-label {
  font-weight: 400;
}

</style>