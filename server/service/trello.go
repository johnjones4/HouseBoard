package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type TrelloConfiguration struct {
	APIKey string
	Token  string
}
type Trello struct {
	configuration TrelloConfiguration
}

func NewTrello(configuration TrelloConfiguration) *Trello {
	return &Trello{
		configuration: configuration,
	}
}

type card struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name"`
	List string `json:"idList"`
}

func (t *Trello) NewCard(listId, name string) (string, error) {
	cardBytes, err := json.Marshal(card{
		Name: name,
		List: listId,
	})
	log.Println(string(cardBytes))
	if err != nil {
		return "", err
	}

	u := url.URL{
		Scheme: "https",
		Host:   "api.trello.com",
		Path:   "/1/cards",
		RawQuery: url.Values{
			"key":   {t.configuration.APIKey},
			"token": {t.configuration.Token},
		}.Encode(),
	}

	res, err := http.Post(u.String(), "application/json", io.NopCloser(bytes.NewBuffer(cardBytes)))
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", errors.New(string(body))
	}

	var c card
	err = json.Unmarshal(body, &c)
	if err != nil {
		return "", err
	}

	cardUrl := fmt.Sprintf("https://trello.com/c/%s", c.Id)

	return cardUrl, nil
}
