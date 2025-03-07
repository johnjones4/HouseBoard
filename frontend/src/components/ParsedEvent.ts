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

export const sortEvents = (events: ParsedEvent[]): ParsedEvent[] => {
  return events.sort((a: ParsedEvent, b: ParsedEvent): number => {
    const aWeight = a.start === a.end ? Number.MAX_SAFE_INTEGER : a.parsedStart.getTime();
    const bWeight = b.start === b.end ? Number.MAX_SAFE_INTEGER : b.parsedStart.getTime();
    return aWeight - bWeight;
  });
}