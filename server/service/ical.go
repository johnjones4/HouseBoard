package service

import (
	"context"
	"log"
	"main/core"
	"net/http"
	"strings"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/mmcdole/gofeed"
)

type event struct {
	Title string    `json:"title"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Label string    `json:"label"`
}

type iCalConfiguration struct {
	URLs map[string]string `json:"urls"`
}

func (c iCalConfiguration) Empty() bool {
	return len(c.URLs) == 0
}

func (c iCalConfiguration) Service() core.Service {
	return &iCal{c}
}

type iCal struct {
	configuration iCalConfiguration
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
		var events []event
		if iical, ok := cal.(*ical.Calendar); ok {
			events, err = extractIcalEvents(name, iical)
			if err != nil {
				return nil, err
			}
		} else if feed, ok := cal.(*gofeed.Feed); ok {
			events, err = extractRssEvents(name, feed)
			if err != nil {
				return nil, err
			}
		}
		output = append(output, events...)
	}
	return output, nil
}

func (i *iCal) NeedsRefresh() bool {
	return true
}

func getCalendar(url string) (any, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if strings.Index(res.Header.Get("Content-type"), "application/rss") == 0 {
		fp := gofeed.NewParser()
		return fp.Parse(res.Body)
	} else {
		return ical.ParseCalendar(res.Body)
	}
}

func extractIcalEvents(label string, cal *ical.Calendar) ([]event, error) {
	now := time.Now()
	plus60Days := now.Add(time.Hour * 24 * 60)
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
					if (start.After(now) || end.After(now)) && (start.Before(plus60Days) || end.Before(plus60Days)) {
						titleOb := icalEvent.GetProperty(ical.ComponentPropertySummary)
						var title string
						if titleOb != nil {
							title = titleOb.Value
						}
						event := event{
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

func extractRssEvents(label string, feed *gofeed.Feed) ([]event, error) {
	now := time.Now()
	plus60Days := now.Add(time.Hour * 24 * 60)
	events := make([]event, 0)
	for _, item := range feed.Items {
		//Mon, 30 Sep 2024 23:59:59 -0400
		publishedParsed, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.Published)
		if err != nil {
			log.Println(err)
			continue
		}
		if publishedParsed.After(now) && publishedParsed.Before(plus60Days) {
			realdate := time.Date(publishedParsed.Year(), publishedParsed.Month(), publishedParsed.Day(), 0, 0, 0, 0, publishedParsed.Location())
			event := event{
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
