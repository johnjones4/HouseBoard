<template>
  <Tile :position="props.position" name="summary">
    <template v-if="infoStore.info.summary">
      <div v-html="infoStore.info.summary.summary"></div>
    </template>
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

.tile.summary .tile-body {
  font-weight: 300;
  font-size: 50px;
  display: flex;
  justify-content: center;
  flex-direction: column;
}

</style>