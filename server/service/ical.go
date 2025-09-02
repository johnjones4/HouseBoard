package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/mmcdole/gofeed"
)

type Event struct {
	Title string    `json:"title"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Label string    `json:"label"`
}

type iCalCalendar struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type iCalConfiguration struct {
	URLs map[string]iCalCalendar `json:"urls"`
}

func (c iCalConfiguration) Empty() bool {
	return len(c.URLs) == 0
}

func (c iCalConfiguration) Service() *ICal {
	return &ICal{c, nil}
}

type ICal struct {
	configuration iCalConfiguration
	Events        []Event
}

func (i *ICal) Name() string {
	return "ical"
}

func (i *ICal) Refresh(c context.Context) error {
	output := make([]Event, 0)
	for name, info := range i.configuration.URLs {
		cal, err := getCalendar(info.URL)
		if err != nil {
			return err
		}
		var events []Event
		if iical, ok := cal.(*ical.Calendar); ok {
			events, err = extractIcalEvents(name, iical)
			if err != nil {
				return err
			}
		} else if feed, ok := cal.(*gofeed.Feed); ok {
			events, err = extractRssEvents(name, feed)
			if err != nil {
				return err
			}
		}
		output = append(output, events...)
	}
	i.Events = output
	return nil
}

func (i *ICal) NeedsRefresh() bool {
	return i != nil
}

func getCalendar(url string) (any, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if strings.Index(res.Header.Get("Content-type"), "application/rss") == 0 || strings.Index(res.Header.Get("Content-type"), "text/xml") == 0 {
		fp := gofeed.NewParser()
		return fp.Parse(res.Body)
	} else {
		return ical.ParseCalendar(res.Body)
	}
}

func extractIcalEvents(label string, cal *ical.Calendar) ([]Event, error) {
	now := time.Now()
	plus60Days := now.Add(time.Hour * 24 * 60)
	events := make([]Event, 0)
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
					if (start.After(now) || end.After(now)) && (start.Before(plus60Days) || end.Before(plus60Days)) {
						titleOb := icalEvent.GetProperty(ical.ComponentPropertySummary)
						var title string
						if titleOb != nil {
							title = titleOb.Value
						}
						event := Event{
							Title: title,
							Start: start,
							End:   end,
							Label: label,
						}
						events = append(events, event)
					}
				}
			}
		}
	}
	return events, nil
}

var rssDateFormats = []string{
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 MST",
}

func extractRssEvents(label string, feed *gofeed.Feed) ([]Event, error) {
	now := time.Now()
	plus60Days := now.Add(time.Hour * 24 * 60)
	events := make([]Event, 0)
	for _, item := range feed.Items {
		var publishedParsed time.Time
		for _, format := range rssDateFormats {
			var err error
			publishedParsed, err = time.Parse(format, item.Published)
			if err == nil {
				break
			}
		}
		if publishedParsed.IsZero() {
			log.Printf("cannot parse date: %s", item.Published)
			continue
		}
		if publishedParsed.After(now) && publishedParsed.Before(plus60Days) {
			realdate := time.Date(publishedParsed.Year(), publishedParsed.Month(), publishedParsed.Day(), 0, 0, 0, 0, publishedParsed.Location())
			event := Event{
				Title: item.Title,
				Start: realdate,
				End:   realdate,
				Label: label,
			}
			events = append(events, event)
		}
	}
	return events, nil
}

func (i *ICal) StateForPrompt() *string {
	if i == nil {
		return nil
	}

	if len(i.Events) == 0 {
		return nil
	}

	now := time.Now()
	horizon := now.Add(time.Hour * 72)

	calendarToEvents := make(map[string][]Event)

	for _, event := range i.Events {
		if event.Start.Before(horizon) && event.Start.After(now) {
			array, ok := calendarToEvents[event.Label]
			if !ok {
				array = make([]Event, 0, 1)
			}
			calendarToEvents[event.Label] = append(array, event)
		}
	}

	if len(calendarToEvents) == 0 {
		return nil
	}

	var calendarSummary strings.Builder

	calendarSummary.WriteString("Calendar of Events:\n")

	for label, events := range calendarToEvents {
		calendarInfo, ok := i.configuration.URLs[label]
		if !ok {
			continue
		}
		calendarSummary.WriteString(fmt.Sprintf("- Calendar Name: %s\n", label))
		calendarSummary.WriteString(fmt.Sprintf("- Calendar Description: %s\n", calendarInfo.Description))
		calendarSummary.WriteString("- Events:\n")
		for i, event := range events {
			calendarSummary.WriteString(fmt.Sprintf("  %d: %s at %s\n", i+1, event.Title, event.Start.In(loc).String()))
		}
	}

	str := calendarSummary.String()

	return &str
}
