<template>
  <Tile :position="props.position">
  </Tile>
</template>

<script lang="ts" setup>

import type { Position } from '../Position';
import Tile from '../Tile.vue';

const props = defineProps<{
  position: Position;
}>();

</script>