<template>
  <Tile :position="props.position" name="radar">
    <div 
      v-if="infoStore.info.forecast && infoStore.info.forecast.radarURL"
      class="radar-display"
      :style="{
        backgroundImage: `url(${infoStore.info.forecast.radarURL}?t=${(new Date().getTime() / 10000).toFixed(0)})`
      }"
    ></div>
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