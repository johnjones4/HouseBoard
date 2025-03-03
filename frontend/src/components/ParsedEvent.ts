import type { components } from "../stores/swagger";

type Event = components['schemas']['Event'];

export interface ParsedEvent extends Event {
  parsedStart: Date;
  parsedEnd: Date;
}

export const mapToParsedEvents = (events: Event[]): ParsedEvent[] => {
  return events.map(e => ({
    ...e,
    parsedStart: new Date(Date.parse(e.start)),
    parsedEnd: new Date(Date.parse(e.end)),
  }));
}