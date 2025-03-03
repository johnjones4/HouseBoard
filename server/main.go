package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"main/api"
	"main/service"
	"net/http"
	"os"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	configBytes, err := os.ReadFile(os.Getenv("CONFIG_FILE"))
	if err != nil {
		panic(err)
	}

	var cfg service.Configuration
	err = json.Unmarshal(configBytes, &cfg)
	if err != nil {
		panic(err)
	}

	services := cfg.Services()

	ctx, cancel := context.WithCancel(context.Background())
	go services.Start(ctx, log)

	r := api.New(log, services)
	err = http.ListenAndServe(os.Getenv("HTTP_HOST"), r)
	if err != nil {
		slog.Error("error listening", slog.Any("error", err))
	}
	cancel()

}
