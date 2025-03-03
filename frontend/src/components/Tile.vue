<template>
  <div 
    :class="['tile', props.name ? props.name : false].filter(c => !!c).join(' ')"
    :style="{
      gridColumnStart: props.position.col,
      gridRowStart: props.position.row,
      gridColumnEnd: props.position.col + props.position.width,
      gridRowEnd: props.position.row + props.position.height,
    }"
  >
    <div v-if="props.title" class="tile-head">
      {{ props.title }}
    </div>
    <div class="tile-body">
      <slot>Tile</slot>
    </div>
  </div>
</template>

<script lang="ts" setup>
import type { Position } from './Position';


const props = defineProps<{
  position: Position;
  name?: string;
  title?: string;
}>();

</script>

<style scoped>

.tile {
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  overflow: hidden;
}

.tile-body {
  flex-grow: 1;
  overflow: hidden;
}

.tile-head {
  font-size: 1.25em;
  padding-bottom: var(--thin-padding);
}

</style>