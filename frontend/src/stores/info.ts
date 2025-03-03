import { defineStore } from 'pinia';
import type { components, paths } from './swagger';
import createClient from 'openapi-fetch';
import { AllColors, type Color } from '../components/colors';

type Info = components['schemas']['Info'];

const client = createClient<paths>({ baseUrl: '/' });

interface InfoStoreState {
  info: Info;
  error?: any;
  eventLabelsToColors: Record<string, Color>
}

export const useInfoStore = defineStore('infoStore', {
  state: (): InfoStoreState => ({
    info: {},
    eventLabelsToColors: {},
  }),
  actions: {
    updateExtras() {
      const labelsSet = new Set<string>();
      if (this.info && this.info.events) {
        this.info.events.events.forEach(e => labelsSet.add(e.label));
        const sorted = Array.from(labelsSet);
        sorted.sort();
        sorted.forEach((label, i) => {
          this.eventLabelsToColors[label] = AllColors[i % AllColors.length];
        });
      }
    },
    async load(): Promise<void> {
      try {
        const {
          data, 
          error,
        } = await client.GET('/info', {});
        if (error) {
          throw error
        }
        this.info = data;
        this.error = undefined;
        this.updateExtras();
      } catch (err) {
        console.error(err);
        this.error = err;
      }
    },
    async start(): Promise<void> {
      setInterval(async () => {
        await this.load();
      }, 1000 * 30);
      await this.load();
    }
  }
});