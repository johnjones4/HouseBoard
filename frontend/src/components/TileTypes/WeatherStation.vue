<template>
  <Tile :position="props.position" name="weather-station">
    <div v-if="infoStore.info.weatherStation" class="weather">
      <div class="weather-inner">
        <div class="weather-temp">
          {{ infoStore.info.weatherStation.temperature.toFixed(1) }}&deg;
        </div>
        <div class="weather-details">
          <p>
            <strong>Pressure:</strong>
            {{ infoStore.info.weatherStation.pressure.toFixed(2) }}inHg
          </p>
          <p>
            <strong>Humidity:</strong>
            {{ infoStore.info.weatherStation.relativeHumidity.toFixed(2) }}%
          </p>
        </div>
      </div>
      <div 
        class="weather-wind"
        :style="{
          transform: `rotate(${infoStore.info.weatherStation.vaneDirection}deg)`
        }"
      >
        {{ infoStore.info.weatherStation.anemometerAverage.toFixed(1) }}mph
      </div>
    </div>
  </Tile>
</template>

<script lang="ts" setup>

import { useInfoStore } from '../../stores/info';
import type { Position } from '../Position';
import Tile from '../Tile.vue';

const infoStore = useInfoStore();

const props = defineProps<{
  position: Position;
}>();


</script>

<style>

.tile.weather-station .tile-body {
  display: flex;
  justify-content: center;
  align-items: center;
}

.tile.weather-station .weather {
  aspect-ratio: 1/1;
  height: 80%;
  border: dotted 2px var(--color-text-light);
  border-radius: 50%;
  position: relative;
}

.tile.weather-station .weather .weather-wind {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  transform-origin: center center;
  text-align: center;
  line-height: 3em;
}

.tile.weather-station .weather .weather-wind::before {
  content: '▼';
  position: absolute;
  top: -1px;
  font-size: 30px;
  line-height: 0;
  left: 0;
  right: 0;
  text-align: center;
}

.tile.weather-station .weather .weather-inner {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
}

.tile.weather-station .weather .weather-inner .weather-temp {
  font-size: 100px;
  text-align: center;
}

.tile.weather-station .weather .weather-inner .weather-details {
  text-align: center;
}

.tile.weather-station .weather .weather-inner .weather-details p {
  margin: 0;
  padding: 0.5em 0 0 0;
}

</style>