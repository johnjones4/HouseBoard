package service

import (
	"context"
	"encoding/json"
	"io"
	"main/core"
	"net/http"
	"net/url"
)

type trelloConfiguration struct {
	APIKey string   `json:"apiKey"`
	Token  string   `json:"token"`
	Lists  []string `json:"lists"`
}

func (c trelloConfiguration) Empty() bool {
	return c.APIKey == "" && c.Token == "" && len(c.Lists) == 0
}

func (c trelloConfiguration) Service() core.Service {
	return &trello{c}
}

type trello struct {
	configuration trelloConfiguration
}

type card struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type list struct {
	Cards []card `json:"cards"`
	Name  string `json:"name"`
}

func (t *trello) Name() string {
	return "trello"
}

func (t *trello) Info(c context.Context) (interface{}, error) {
	resp := make([]list, len(t.configuration.Lists))
	for i, id := range t.configuration.Lists {
		list, err := t.getList(id)
		if err != nil {
			return nil, err
		}

		cards, err := t.getCards(id)
		if err != nil {
			return nil, err
		}

		list.Cards = cards

		resp[i] = list
	}
	return resp, nil
}

func (t *trello) getList(id string) (list, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "api.trello.com",
		Path:   "/1/lists/" + id,
		RawQuery: url.Values{
			"key":   {t.configuration.APIKey},
			"token": {t.configuration.Token},
		}.Encode(),
	}

	res, err := http.Get(u.String())
	if err != nil {
		return list{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return list{}, err
	}

	var l list
	err = json.Unmarshal(body, &l)
	if err != nil {
		return list{}, err
	}

	return l, nil
}

func (t *trello) getCards(id string) ([]card, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "api.trello.com",
		Path:   "/1/lists/" + id + "/cards",
		RawQuery: url.Values{
			"key":   {t.configuration.APIKey},
			"token": {t.configuration.Token},
		}.Encode(),
	}

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var cardsResp []card
	err = json.Unmarshal(body, &cardsResp)
	if err != nil {
		return nil, err
	}

	return cardsResp, nil
}

func (t *trello) NeedsRefresh() bool {
	return true
}
