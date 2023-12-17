package main

import (
	"encoding/json"
	"main/api"
	"main/core"
	"main/service"
	"net/http"
	"os"
	"strconv"

	"github.com/johnjones4/errorbot"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	chatId, err := strconv.Atoi(os.Getenv("TELEGRAM_CHAT_ID"))
	if err != nil {
		panic(err)
	}
	bot := errorbot.New(
		"houseboard",
		os.Getenv("TELEGRAM_TOKEN"),
		chatId,
	)

	config := zap.NewDevelopmentConfig()
	l, err := config.Build(zap.Hooks(bot.ZapHook([]zapcore.Level{
		zapcore.FatalLevel,
		zapcore.PanicLevel,
		zapcore.DPanicLevel,
		zapcore.ErrorLevel,
		zapcore.WarnLevel,
	})))
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
