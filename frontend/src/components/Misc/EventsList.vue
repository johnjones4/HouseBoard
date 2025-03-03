<template>
  <div class="events-list">
    <div
      v-for="e, j in props.events"
      :key="j"
      class="events-list-event"
    >
      <span class="events-list-bullet" :style="{backgroundColor: infoStore.eventLabelsToColors[e.label]}"></span>
      {{ timeString(e) }} {{ e.title }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { useInfoStore } from '../../stores/info';
import { hoursMinutesString } from '../../util';
import type { ParsedEvent } from '../ParsedEvent';

const infoStore = useInfoStore();

const props = defineProps<{
  events: ParsedEvent[];
}>();

const timeString = (event: ParsedEvent): string => {
  const start = hoursMinutesString(event.parsedStart)
  const end = hoursMinutesString(event.parsedEnd)
  return start === end ? '' : `${start} | `
}

</script>

<style scoped>

.events-list-bullet {
  display: inline-block;
  width: 0.5em;
  height: 0.5em;
  border-radius: 50%;
}


</style>