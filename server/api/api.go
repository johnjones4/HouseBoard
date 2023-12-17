package api

import (
	"context"
	"encoding/json"
	"log"
	"main/core"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func handleError(w http.ResponseWriter, err error, status int) {
	log.Println(err)
	http.Error(w, http.StatusText(status), status)
}

func jsonResponse(w http.ResponseWriter, j interface{}) {
	bytes, err := json.Marshal(j)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.Write(bytes)
}

const defaultWait = time.Second * 5

func New(services []core.Service, log *zap.SugaredLogger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	output := make(map[string]interface{})
	outputLock := new(sync.RWMutex)
	go func() {
		log.Info("Starting poller")
		for {
			for _, service := range services {
				name := service.Name()
				if _, ok := output[name]; !ok || service.NeedsRefresh() {
					var err error
					waitDelay := defaultWait
					for err != nil || waitDelay == defaultWait {
						log.Infof("Updating %s", name)
						waitDelay = waitDelay * 2
						info, err := service.Info(context.Background())
						if err != nil {
							log.Errorf("Error: \"%s\". Sleeping %s", err.Error(), waitDelay.String())
							time.Sleep(waitDelay)
							if waitDelay > time.Second*30 {
								break
							}
						} else {
							outputLock.Lock()
							output[name] = info
							outputLock.Unlock()
						}
					}
				}
			}
			time.Sleep(time.Minute * 5)
		}
	}()

	r.Get("/api/info", func(w http.ResponseWriter, r *http.Request) {
		outputLock.RLock()
		jsonResponse(w, output)
		outputLock.RUnlock()
	})
	return r
}
