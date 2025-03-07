package service

import (
	"bufio"
	"context"
	"errors"
	"os"
	"strings"
	"time"
)

var baseDate = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

type TriviaConfiguration struct {
	Files []string `json:"files"`
}

func (c TriviaConfiguration) Empty() bool {
	return len(c.Files) == 0
}

func (c TriviaConfiguration) Service() *Trivia {
	return &Trivia{
		TriviaConfiguration: c,
	}
}

type Trivia struct {
	TriviaConfiguration

	lastRefresh time.Time
	Current     *TriviaQuestion
	Previous    *TriviaQuestion
}

func (f *Trivia) Name() string {
	return "trivia"
}

func (f *Trivia) Refresh(c context.Context) error {
	now := time.Now()
	current, err := f.getTriviaQuestion(now)
	if err != nil {
		return err
	}
	f.Current = current
	previous, err := f.getTriviaQuestion(now.AddDate(0, 0, -1))
	if err != nil {
		return err
	}
	f.Previous = previous
	f.lastRefresh = now
	return nil
}

func (f *Trivia) NeedsRefresh() bool {
	now := time.Now()
	return now.Format(time.DateOnly) != f.lastRefresh.Format(time.DateOnly)
}

func (f *Trivia) StateForPrompt() *string {
	return nil
}

type TriviaQuestion struct {
	Question string
	Choices  []string
	Answer   int
}

var errorRanOutOfFile = errors.New("ran out of file to scan")

func (f *Trivia) getTriviaQuestion(date time.Time) (*TriviaQuestion, error) {
	daysSinceStart := int(date.Sub(baseDate) / (time.Hour * 24))
	nFile := daysSinceStart % len(f.Files)
	nEntry := int(float64(daysSinceStart) / float64(len(f.Files)))

	file, err := os.Open(f.Files[nFile])
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	currentQuestion := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Index(line, "#Q ") == 0 {
			if currentQuestion == nEntry {
				question := line[3:]
				if !scanner.Scan() {
					return nil, errorRanOutOfFile
				}
				correctAnswer := strings.TrimSpace(scanner.Text()[2:])
				correctAnserIndex := -1
				choices := make([]string, 0, 4)
			choices:
				for scanner.Scan() {
					line := scanner.Text()
					if strings.TrimSpace(line) == "" {
						break choices
					}
					choice := strings.TrimSpace(line[2:])
					if choice == correctAnswer {
						correctAnserIndex = len(choices)
					}
					choices = append(choices, choice)
				}
				return &TriviaQuestion{
					Question: question,
					Answer:   correctAnserIndex,
					Choices:  choices,
				}, nil
			} else {
				currentQuestion++
			}
		}
	}

	return nil, errorRanOutOfFile
}
