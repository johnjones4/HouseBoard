<template>
  <Tile :position="props.position" name="trivia">
    <template v-if="infoStore.info.trivia">
      <div class="trivia-section trivia-current">
        <div class="trivia-label">Question:</div>
        <div class="trivia-question">{{ infoStore.info.trivia.question }}</div>
        <div class="trivia-label">Choices:</div>
        <div class="trivia-choices">
          <ol>
            <li v-for="c in infoStore.info.trivia.choices" :key="c">{{ c }}</li>
          </ol>
        </div>
      </div>
      <div class="trivia-section trivia-previous">
        <div class="trivia-label">Yesterday's Question:</div>
        <div class="trivia-question">{{ infoStore.info.trivia.previousQuestion }}</div>
        <div class="trivia-label">Answer:</div>
        <div class="trivia-answer">{{ infoStore.info.trivia.previousAnswer }}</div>
      </div>
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

.tile.trivia .tile-body {
  font-weight: 300;
  font-size: 60px;
  display: flex;
  justify-content: center;
  flex-direction: column;
}

.tile.trivia .trivia-label {
  font-weight: 500;
  font-size: 0.8em;
  margin-top: 1em;
}

.tile.trivia ol {
  margin: 0.25em 0 0 0;
  padding: 0 0 0 1.5em;
}

.tile.trivia ol li {
  margin: 0;
  list-style: upper-alpha;
}

.tile.trivia .trivia-current {
  font-size: 60px;
}

.tile.trivia .trivia-previous {
  font-size: 30px;
  margin-top: 1em;
}

.tile.trivia .trivia-section {
  max-width: 1500px;
  width: 100%;
  margin-left: auto;
  margin-right: auto;
}

</style>