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
          v-else-if="tile.tileType === TileType.radar" 
          :position="tile.position"
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
      </template>
    </div>
    <div v-if="infoStore && infoStore.info && infoStore.info.weatherStation" class="frame-footer">
      <div class="frame-footer-prop">
        <div class="frame-footer-prop-value" ref="dateTimeEl">
          {{ new Date().toLocaleString() }}
        </div>
        <div class="frame-footer-prop-label">
          Date/Time
        </div>
      </div>
      <div class="frame-footer-prop">
        <div class="frame-footer-prop-value">
          {{ infoStore.info.weatherStation.temperature.toFixed(2) }}&deg;
        </div>
        <div class="frame-footer-prop-label">
          Temperature
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
import Summary from './TileTypes/Summary.vue';
import Traffic from './TileTypes/Traffic.vue';
import Trello from './TileTypes/Trello.vue';
import WeatherStation from './TileTypes/WeatherStation.vue';
import type { FrameProps } from './FrameProps';
import Agenda from './TileTypes/Agenda.vue';
import { useInfoStore } from '../stores/info';
import { ref } from 'vue';

const props = defineProps<FrameProps>();

const dateTimeEl = ref<null|HTMLDivElement>(null);

const infoStore = useInfoStore();

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
}

.frame-footer-prop:first-child {
  margin-left: var(--thin-padding);
  border: none;
  padding: 0;
}

.frame-footer-prop-value {
  font-size: 2em;
  font-optical-sizing: auto;
  font-family: "Doto", monospace;
  font-weight: bold;
}

.frame-footer-prop-label {
  font-weight: 400;
}

</style>