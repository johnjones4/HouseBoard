package service

import (
	"context"
	"main/core"
)

type fileConfiguration struct {
	Files []string `json:"files"`
}

func (c fileConfiguration) Empty() bool {
	return len(c.Files) == 0
}

func (c fileConfiguration) Service() core.Service {
	return &file{
		fileConfiguration: c,
	}
}

type file struct {
	fileConfiguration
}

func (f *file) Name() string {
	return "file"
}

func (f *file) Info(c context.Context) (interface{}, error) {
	return f.fileConfiguration, nil
}

func (f *file) NeedsRefresh() bool {
	return false
}
