<template>
  <Tile :position="props.position" name="calendar">
    <div v-if="infoStore.info.events && infoStore.info.events.events" class="calendar-calendar">
      <div
        v-for="d in daysOfWeek"
        :key="d"
        class="calendar-header"
      >
        {{ d }}
      </div>
      <template
        v-for="event, i in calendar"
        :key="event.date || i"
      >
        <div 
          v-if="event.date"
          class="calendar-item"
        >
          <div class="calendar-item-date">
            {{ event.date.getDate() }}
          </div>
          <EventsList :events="event.events" />
        </div>
        <div 
          v-else 
          class="calendar-item calendar-item-empty">
        </div>
      </template>
    </div>
  </Tile>
</template>

<script lang="ts" setup>

import { computed } from 'vue';
import { useInfoStore } from '../../stores/info';
import type { Position } from '../Position';
import Tile from '../Tile.vue';
import { mapToParsedEvents, type ParsedEvent } from '../ParsedEvent';
import EventsList from '../Misc/EventsList.vue';


const infoStore = useInfoStore();

const props = defineProps<{
  position: Position;
}>();

const daysOfWeek = ['Su', 'M', 'Tu', 'W', 'Th', 'F', 'Sa'];

interface CalendarItem {
  events: ParsedEvent[]
  date: Date | null
}

const calendar = computed((): CalendarItem[] => {
  if (!infoStore.info.events) {
    return [];
  }
  const now = new Date();
  let curDay = now;
  const array = [] as CalendarItem[];
  for (let row = 0; row < 5; row++) {
    if (curDay.getDay() > 0) {
      for (let i = 0; i < curDay.getDay(); i++) {
        array.push({
          events: [],
          date: null,
        })
      }
    }
    const parsedEvents = mapToParsedEvents(infoStore.info.events.events);
    for (let col = curDay.getDay(); col < 7; col++) {      
      const date = curDay;
      let events = parsedEvents.filter(event => eventOccursOnDay(event, date));
      events = events.sort((a: ParsedEvent, b: ParsedEvent): number => {
        const aWeight = a.start === a.end ? Number.MAX_SAFE_INTEGER : a.parsedStart.getTime();
        const bWeight = b.start === b.end ? Number.MAX_SAFE_INTEGER : b.parsedStart.getTime();
        return aWeight - bWeight;
      });
      array.push({
        events,
        date,
      });
      curDay = new Date(curDay.getTime() + (1000*60*60*24));
    }
  }
  return array;
})


const eventOccursOnDay = (event: ParsedEvent, day: Date): boolean => {
  const startOfDay = new Date(day.getFullYear(), day.getMonth(), day.getDate(), 0, 0, 0, 0)
  const endOfDay = new Date(startOfDay.getTime() + (1000*60*60*24))
  return (event.parsedStart.getTime() >= startOfDay.getTime() && event.parsedStart.getTime() < endOfDay.getTime())
    || (event.parsedEnd.getTime() > startOfDay.getTime() && event.parsedEnd.getTime() <= endOfDay.getTime())
    || (startOfDay.getTime() >= event.parsedStart.getTime() && endOfDay.getTime() <= event.parsedEnd.getTime())
}

</script>

<style>

.tile.calendar .tile-body {
  display: flex;
  font-size: 12px;
}

.tile.calendar .calendar-calendar {
  flex-grow: 1;
  display: grid;
  grid-template-columns: repeat(7, minmax(0, 1fr));
  grid-template-rows: 2em repeat(5, minmax(0, 1fr));
}

.tile.calendar .calendar-header {
  font-weight: bold;
  padding-left: var(--thin-padding);
}

.tile.calendar .calendar-item {
  padding: var(--thin-padding);
  border-left: solid 1px var(--color-text-light);
  border-top: solid 1px var(--color-text-light);
  overflow: hidden;
}

.tile.calendar .calendar-item:nth-child(7n) {
  border-right: solid 1px var(--color-text-light);
}

.tile.calendar .calendar-item:nth-last-child(-n+7) {
  border-bottom: solid 1px var(--color-text-light);
}

.tile.calendar .calendar-item.calendar-item-empty {
  background: none;
}

.tile.calendar .calendar-item-date {
  text-align: right;
  font-size: 1.25em;
  padding-bottom: var(--thin-padding);
}

</style>