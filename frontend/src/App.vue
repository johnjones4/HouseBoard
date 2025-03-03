<template>
  <Frame 
    :title="tileSets[currentIndex].title"
    :tiles="tileSets[currentIndex].tiles"
  />
  <div
    class="countdown"
    ref="countdown"
  ></div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useInfoStore } from './stores/info';
import Frame from './components/Frame.vue';
import { TileType } from './components/TileType';
import type { FrameProps } from './components/FrameProps';

const infoStore = useInfoStore();

const tileSets = [
  {
    title: 'Events',
    tiles: [
      {
        tileType: TileType.agenda,
        position: {
          row: 1,
          col: 1,
          width: 2,
          height: 6,
        }
      },
      {
        tileType: TileType.calendar,
        position: {
          row: 1,
          col: 3,
          width: 10,
          height: 6,
        }
      },
    ]
  },
  {
    title: 'Current Conditions',
    tiles: [
      {
        tileType: TileType.files,
        position: {
          row: 1,
          col: 1,
          width: 8,
          height: 3,
        }
      },
      {
        tileType: TileType.radar,
        position: {
          row: 1,
          col: 9,
          width: 4,
          height: 3,
        }
      },
      {
        tileType: TileType.weatherStation,
        position: {
          row: 4,
          col: 1,
          width: 3,
          height: 3,
        }
      },
      {
        tileType: TileType.forecast,
        position: {
          row: 4,
          col: 4,
          width: 9,
          height: 3,
        }
      },
    ]
  },
  {
    title: 'Summary',
    tiles: [
      {
        tileType: TileType.summary,
        position: {
          row: 1,
          col: 1,
          width: 12,
          height: 6,
        },
      },
    ],
  },
] as FrameProps[];
const currentIndex = ref(0);
const nextChange = ref(0);
const countdown = ref<null|HTMLDivElement>(null);
const delay = 20000;

const tick = () => {
  const now = new Date().getTime();
  if (nextChange.value <= now) {
    nextChange.value = 0;
    currentIndex.value = (currentIndex.value + 1) % tileSets.length;
  }
  if (nextChange.value === 0) {
    nextChange.value = now + delay;
  }
  if (countdown.value) {
    countdown.value.style.width = `${(nextChange.value - now) / delay * 100}%`;
  }
  requestAnimationFrame(tick);
}

onMounted(() => {
  infoStore.start();
});

requestAnimationFrame(tick);
</script>

<style scoped>

.countdown {
  position: fixed;
  z-index: 10000;
  bottom: 0;
  right: 0;
  height: var(--progress-line-thickness);
  background-color: var(--color-red);
}

</style>
