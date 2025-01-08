package main

import (
	"encoding/json"
	"main/api"
	"main/core"
	"main/service"
	"net/http"
	"os"

	"go.uber.org/zap"
)

func main() {
	config := zap.NewDevelopmentConfig()
	l, err := config.Build()
	if err != nil {
		panic(err)
	}

	defer l.Sync()
	log := l.Sugar()

	configBytes, err := os.ReadFile(os.Getenv("CONFIG_FILE"))
	if err != nil {
		panic(err)
	}

	var cfg service.Configuration
	err = json.Unmarshal(configBytes, &cfg)
	if err != nil {
		panic(err)
	}

	services := make([]core.Service, 0)
	for _, config := range cfg.Configurations() {
		if !config.Empty() {
			services = append(services, config.Service())
		}
	}

	router := api.New(services, log)
	err = http.ListenAndServe(os.Getenv("HTTP_HOST"), router)
	panic(err)
}
