<template>
  <Tile :position="props.position">
  </Tile>
</template>

<script lang="ts" setup>

import { storeToRefs } from 'pinia';
import { useInfoStore } from '../../stores/info';
import type { Position } from '../Position';
import Tile from '../Tile.vue';
import { onMounted, watch } from 'vue';

const infoStore = useInfoStore();
const { info } = storeToRefs(infoStore)

const props = defineProps<{
  position: Position;
}>();

const update = () => {
  
}

onMounted(() => update());
watch(info, () => update());

</script>