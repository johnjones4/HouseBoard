package main

import (
	"encoding/json"
	"main/api"
	"main/core"
	"main/service"
	"net/http"
	"os"
)

func main() {
	configBytes, err := os.ReadFile(os.Getenv("CONFIG_FILE"))
	if err != nil {
		panic(err)
	}

	var config service.Configuration
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}

	services := make([]core.Service, 0)
	for _, config := range config.Configurations() {
		if !config.Empty() {
			services = append(services, config.Service())
		}
	}

	router := api.New(services)
	err = http.ListenAndServe(os.Getenv("HTTP_HOST"), router)
	panic(err)
}
