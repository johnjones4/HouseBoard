package service

import (
	"context"
)

type FileConfiguration struct {
	Files []string `json:"files"`
}

func (c FileConfiguration) Empty() bool {
	return len(c.Files) == 0
}

func (c FileConfiguration) Service() *Files {
	return &Files{
		Files: c.Files,
	}
}

type Files struct {
	Files []string
}

func (f *Files) Name() string {
	return "file"
}

func (f *Files) Refresh(c context.Context) error {
	return nil
}

func (f *Files) NeedsRefresh() bool {
	return false
}

func (f *Files) StateForPrompt() *string {
	return nil
}
