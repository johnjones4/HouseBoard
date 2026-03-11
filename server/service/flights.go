package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"main/core"
	"net/http"
	"net/url"
	"time"
)

type flightsConfiguration struct {
	ClientId     string          `json:"clientId"`
	ClientSecret string          `json:"clientSecret"`
	Min          core.Coordinate `json:"min"`
	Max          core.Coordinate `json:"max"`
}

func (c flightsConfiguration) Empty() bool {
	return c.ClientId == "" || c.ClientSecret == ""
}

func (c flightsConfiguration) Service() *Flights {
	return &Flights{
		configuration: c,
	}
}

type Flights struct {
	configuration   flightsConfiguration
	FlightsResponse []flightStatus
	LastUpdated     *time.Time
	accessToken     *string
	tokenExpires    *time.Time
}

func (w *Flights) Name() string {
	return "flights"
}

func (w *Flights) Refresh(c context.Context) error {
	if w.accessToken == nil || w.tokenExpires == nil || w.tokenExpires.Before(time.Now()) {
		reqbody := url.Values{
			"grant_type":    []string{"client_credentials"},
			"client_id":     []string{w.configuration.ClientId},
			"client_secret": []string{w.configuration.ClientSecret},
		}.Encode()
		res, err := http.Post("https://auth.opensky-network.org/auth/realms/opensky-network/protocol/openid-connect/token", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(reqbody)))
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("bad auth response: %d", res.StatusCode)
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		var authResponse flightAuthResponse
		err = json.Unmarshal(body, &authResponse)
		if err != nil {
			return err
		}
		w.accessToken = &authResponse.AccessToken
		w.tokenExpires = new(time.Now().Add(time.Second * time.Duration(authResponse.ExpiresIn/2)))
	}

	params := url.Values{
		"lamax":    []string{fmt.Sprint(w.configuration.Max.Latitude)},
		"lamin":    []string{fmt.Sprint(w.configuration.Min.Latitude)},
		"lomax":    []string{fmt.Sprint(w.configuration.Max.Longitude)},
		"lomin":    []string{fmt.Sprint(w.configuration.Min.Longitude)},
		"extended": []string{"1"},
	}
	req, err := http.NewRequest("GET", "https://opensky-network.org/api/states/all?"+params.Encode(), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+*w.accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response: %d", res.StatusCode)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var info flightsResponse
	err = json.Unmarshal(body, &info)
	if err != nil {
		return err
	}

	if info.States != nil {
		w.FlightsResponse = *info.States
	} else {
		w.FlightsResponse = []flightStatus{}
	}
	w.LastUpdated = new(time.Now())

	return nil
}

func (w *Flights) NeedsRefresh() bool {
	return w.LastUpdated == nil || time.Since(*w.LastUpdated) > time.Second*30
}

func (w *Flights) StateForPrompt() *string {
	return nil
}

type flightStatus []any

func (s flightStatus) GetCallsign() *string {
	if 1 < len(s) {
		return new(s[1].(string))
	}
	return nil
}

func (s flightStatus) GetAlitude() *float64 {
	if 7 < len(s) {
		return new(s[7].(float64))
	}
	return nil
}

type flightsResponse struct {
	States *[]flightStatus `json:"states"`
}

type flightAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
