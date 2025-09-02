package service

import (
	"context"
	"fmt"
	"log/slog"
	"main/core"
	"strings"
	"sync"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/openai/openai-go"

	_ "embed"
)

//go:embed prompt.txt
var promptStart string

var scrubber = bluemonday.UGCPolicy()

type Services struct {
	ICal           *ICal
	NOAA           *NOAA
	OpenWeatherMap *OpenWeatherMap
	Traffic        *Traffic
	Trello         *Trello
	Files          *Files
	WeatherStation *WeatherStation
	Lock           sync.RWMutex
	SunriseSunset  *SunriseSunset
	Trivia         *Trivia

	summary string
}

func (s *Services) All() []core.Service {
	return []core.Service{
		s.ICal,
		s.NOAA,
		s.OpenWeatherMap,
		s.Traffic,
		s.Trello,
		s.Files,
		s.WeatherStation,
		s.SunriseSunset,
		s.Trivia,
	}
}

func (s *Services) Summary() string {
	return s.summary
}

func (s *Services) refresh(ctx context.Context, log *slog.Logger) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	log.Info("refreshing services")
	for _, service := range s.All() {
		if service != nil && service.NeedsRefresh() {
			err := service.Refresh(ctx)
			if err != nil {
				log.Error("error refreshing service", slog.String("service", service.Name()), slog.Any("error", err))
			}
		}
	}
}

func (s *Services) summarize(ctx context.Context, log *slog.Logger) {
	log.Info("summarizing services")

	var prompt strings.Builder
	prompt.WriteString(promptStart)
	prompt.WriteString("\n\n")
	prompt.WriteString(fmt.Sprintf("Current Date: %s", time.Now().In(loc).Format(time.ANSIC)))
	prompt.WriteString("\n\n")
	prompt.WriteString("Current Summary: ")
	if s.summary == "" {
		prompt.WriteString("<NONE>")
	} else {
		prompt.WriteString(s.summary)
	}
	prompt.WriteString("\n\n")
	prompt.WriteString("Current State:")
	for _, service := range s.All() {
		if service != nil {
			state := service.StateForPrompt()
			if state != nil {
				prompt.WriteString("\n\n")
				prompt.WriteString(*state)
			}
		}
	}

	log.Info("prompt generated", slog.String("prompt", prompt.String()))

	client := openai.NewClient()
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt.String()),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	})
	if err != nil || len(chatCompletion.Choices) == 0 {
		log.Error("error summarizing", slog.Any("error", err))
		return
	}

	log.Info("prompt response", slog.String("chatCompletion", chatCompletion.Choices[0].Message.Content))

	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.summary = scrubber.Sanitize(chatCompletion.Choices[0].Message.Content)
}

func (s *Services) Start(ctx context.Context, log *slog.Logger) {
	s.refresh(ctx, log)
	s.summarize(ctx, log)
	refreshTicker := time.Tick(time.Second * 30)
	summarizeTicker := time.Tick(time.Hour)
	for {
		select {
		case <-ctx.Done():
			return
		case <-refreshTicker:
			s.refresh(ctx, log)
		case <-summarizeTicker:
			s.summarize(ctx, log)
		}
	}
}
