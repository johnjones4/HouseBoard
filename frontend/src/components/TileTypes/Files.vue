<template>
  <Tile :position="props.position" name="files">
    <div 
      v-if="infoStore.info && infoStore.info.files"
      v-for="file in infoStore.info.files.files"
      :key="file"
      class="file-display"
      :style="{
        backgroundImage: `url(${file}?t=${(new Date().getTime() / 10000).toFixed(0)})`
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

.tile.files .tile-body {
  display: flex;
}

.tile.files .tile-body .file-display {
  width: 100%;
  height: auto;
  max-width: 100%;
  max-height: 100%;
  background-size: cover;
  background-position: center center;
  background-repeat: no-repeat;
}

</style>