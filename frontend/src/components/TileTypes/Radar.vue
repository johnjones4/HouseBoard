<template>
  <Tile :position="props.position" name="radar">
    <div 
      v-if="inforStore.info && inforStore.info.forecast"
      class="radar-display"
      :style="{
        backgroundImage: `url(${url}?t=${(new Date().getTime() / 10000).toFixed(0)})`
      }"
    ></div>
  </Tile>
</template>

<script lang="ts" setup>

import { computed } from 'vue';
import { useInfoStore } from '../../stores/info';
import type { Position } from '../Position';
import Tile from '../Tile.vue';
import { RadarType } from '../RadarType';

const inforStore = useInfoStore();

const props = defineProps<{
  position: Position;
  radarType: RadarType
}>();

const url = computed((): string => {
  if (!inforStore.info || !inforStore.info.forecast) {
    return '';
  }
  switch (props.radarType) {
    case RadarType.local:
      return inforStore.info.forecast.localRadarURL;
    case RadarType.national:
      return inforStore.info.forecast.nationalRadarURL;
    default:
      return ''
  }
})

</script>

<style>

.tile.radar .tile-body {
  display: flex;
}

.tile.radar .tile-body .radar-display {
  width: 100%;
  height: auto;
  max-width: 100%;
  max-height: 100%;
  background-size: contain;
  background-position: center center;
  background-repeat: no-repeat;
}

</style>