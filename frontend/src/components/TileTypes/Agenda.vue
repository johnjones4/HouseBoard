<template>
  <Tile :position="props.position" name="agenda" title="Today's Agenda">
    <EventsList :events="agenda" />
  </Tile>
</template>

<script lang="ts" setup>

import { computed } from 'vue';
import { useInfoStore } from '../../stores/info';
import type { Position } from '../Position';
import Tile from '../Tile.vue';
import EventsList from '../Misc/EventsList.vue';
import { mapToParsedEvents, sortEvents, type ParsedEvent } from '../ParsedEvent';

const infoStore = useInfoStore();

const props = defineProps<{
  position: Position;
}>();

const agenda = computed((): ParsedEvent[] => {
  if (!infoStore.info.events) {
    return [];
  }
  const now = new Date();
  return sortEvents(mapToParsedEvents(infoStore.info.events.events)
    .filter(e => e.parsedStart.getDate() === now.getDate() && e.parsedStart.getMonth() === now.getMonth() && e.parsedStart.getFullYear() === now.getFullYear()));
})

</script>

<style>

</style>