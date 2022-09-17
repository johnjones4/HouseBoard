package service

import (
	"context"
	"main/core"
	"net/http"
	"time"

	ical "github.com/arran4/golang-ical"
)

type event struct {
	Title string    `json:"title"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Label string    `json:"label"`
}

type ICalConfiguration struct {
	URLs map[string]string `json:"urls"`
}

func (c ICalConfiguration) Empty() bool {
	return len(c.URLs) == 0
}

func (c ICalConfiguration) Service() core.Service {
	return &iCal{c}
}

type iCal struct {
	configuration ICalConfiguration
}

func (i *iCal) Name() string {
	return "ical"
}

func (i *iCal) Info(c context.Context) (interface{}, error) {
	output := make([]event, 0)
	for name, url := range i.configuration.URLs {
		cal, err := getCalendar(url)
		if err != nil {
			return nil, err
		}
		events, err := extractEvents(name, cal)
		if err != nil {
			return nil, err
		}
		output = append(output, events...)
	}
	return output, nil
}

func getCalendar(url string) (*ical.Calendar, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return ical.ParseCalendar(res.Body)
}

func extractEvents(label string, cal *ical.Calendar) ([]event, error) {
	events := make([]event, 0)
	for _, comp := range cal.Components {
		if icalEvent, ok := comp.(*ical.VEvent); ok {
			if icalEvent.GetProperty(ical.ComponentPropertyDtStart) != nil && icalEvent.GetProperty(ical.ComponentPropertyDtEnd) != nil {
				end, err := icalEvent.GetEndAt()
				if err != nil {
					return nil, err
				}
				if end.After(time.Now()) {
					start, err := icalEvent.GetStartAt()
					if err != nil {
						return nil, err
					}
					event := event{
						Title: icalEvent.GetProperty(ical.ComponentPropertySummary).Value,
						Start: start,
						End:   end,
						Label: label,
					}
					events = append(events, event)
				}
			}
		}
	}
	return events, nil
}
