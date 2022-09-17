package api

import (
	"encoding/json"
	"log"
	"main/core"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

func New(services []core.Service) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	r.Get("/api/info", func(w http.ResponseWriter, r *http.Request) {
		//TODO context parsing
		output := make(map[string]interface{})
		for _, service := range services {
			info, err := service.Info(r.Context())
			if err != nil {
				handleError(w, err, http.StatusInternalServerError)
				return
			}
			output[service.Name()] = info
		}
		jsonResponse(w, output)
	})
	return r
}
