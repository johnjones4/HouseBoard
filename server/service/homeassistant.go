package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
)

type HomeAssistantConfiguration struct {
	Token     string   `json:"token"`
	Url       string   `json:"url"`
	EntityIds []string `json:"entityIds"`
}

func (c HomeAssistantConfiguration) Empty() bool {
	return c.Token == "" && c.Url == "" && len(c.EntityIds) == 0
}

func (c HomeAssistantConfiguration) Service() *HomeAssistant {
	return &HomeAssistant{
		HomeAssistantConfiguration: c,
		States:                     make(map[string]haEntityState),
	}
}

type HomeAssistant struct {
	HomeAssistantConfiguration
	States map[string]haEntityState
}

type haEntityState struct {
	EntityId string `json:"entity_id"`
	State    string `json:"state"`
}

func (f *HomeAssistant) Name() string {
	return "homeAssistant"
}

func (f *HomeAssistant) Refresh(c context.Context) error {
	req, err := http.NewRequest("GET", f.Url+"/api/states", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", f.Token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", res.StatusCode)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var entities []haEntityState
	err = json.Unmarshal(bytes, &entities)
	if err != nil {
		return err
	}

	for _, entity := range entities {
		if slices.Contains(f.EntityIds, entity.EntityId) {
			f.States[entity.EntityId] = entity
		}
	}

	return nil
}

func (f *HomeAssistant) NeedsRefresh() bool {
	return true
}

func (f *HomeAssistant) StateForPrompt() *string {
	return nil
}
