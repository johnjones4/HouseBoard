package service

import (
	"context"
	"main/core"
	"time"

	"github.com/nathan-osman/go-sunrise"
)

type SunriseSunsetConfiguration struct {
	Location core.Coordinate
}

func (c SunriseSunsetConfiguration) Empty() bool {
	return c.Location.Latitude == 0 && c.Location.Longitude == 0
}

func (c SunriseSunsetConfiguration) Service() *SunriseSunset {
	return &SunriseSunset{
		SunriseSunsetConfiguration: c,
	}
}

type SunriseSunset struct {
	SunriseSunsetConfiguration

	Rise *time.Time
	Set  *time.Time
}

func (f *SunriseSunset) Name() string {
	return "sunriseSunset"
}

func (f *SunriseSunset) Refresh(c context.Context) error {
	now := time.Now().In(loc)
	rise, set := sunrise.SunriseSunset(
		f.Location.Latitude, f.Location.Longitude,
		now.Year(), now.Month(), now.Day(),
	)
	riseLocal := rise.In(loc)
	setLocal := set.In(loc)
	f.Rise = &riseLocal
	f.Set = &setLocal
	return nil
}

func (f *SunriseSunset) NeedsRefresh() bool {
	return true
}

func (f *SunriseSunset) StateForPrompt() *string {
	return nil
}
