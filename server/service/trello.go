package service

import (
	"context"
	"encoding/json"
	"io"
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

func (c trelloConfiguration) Service() *Trello {
	return &Trello{c, nil}
}

type Trello struct {
	configuration trelloConfiguration
	List          []List
}

type Card struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type List struct {
	Cards []Card `json:"cards"`
	Name  string `json:"name"`
}

func (t *Trello) Name() string {
	return "trello"
}

func (t *Trello) Refresh(c context.Context) error {
	resp := make([]List, len(t.configuration.Lists))
	for i, id := range t.configuration.Lists {
		list, err := t.getList(id)
		if err != nil {
			return err
		}

		cards, err := t.getCards(id)
		if err != nil {
			return err
		}

		list.Cards = cards

		resp[i] = list
	}
	t.List = resp
	return nil
}

func (t *Trello) getList(id string) (List, error) {
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
		return List{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return List{}, err
	}

	var l List
	err = json.Unmarshal(body, &l)
	if err != nil {
		return List{}, err
	}

	return l, nil
}

func (t *Trello) getCards(id string) ([]Card, error) {
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

	var cardsResp []Card
	err = json.Unmarshal(body, &cardsResp)
	if err != nil {
		return nil, err
	}

	return cardsResp, nil
}

func (t *Trello) NeedsRefresh() bool {
	return t != nil
}

func (t *Trello) StateForPrompt() *string {
	return nil
}
