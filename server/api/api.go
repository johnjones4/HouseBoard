package api

import (
	"log/slog"
	"main/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	slogchi "github.com/samber/slog-chi"
)

func New(log *slog.Logger, svc *service.Services) *chi.Mux {
	r := chi.NewRouter()

	logConfig := slogchi.Config{
		DefaultLevel:       slog.LevelInfo,
		ClientErrorLevel:   slog.LevelError,
		ServerErrorLevel:   slog.LevelError,
		WithUserAgent:      true,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         false,
		WithTraceID:        false,
		Filters:            nil,
	}

	r.Use(
		middleware.RealIP,
		slogchi.NewWithConfig(log, logConfig),
		middleware.RequestID,
		middleware.Recoverer,
		middleware.Compress(5),
	)

	ctrl := &controller{
		services: svc,
	}
	options := ChiServerOptions{
		BaseRouter: r,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Error("handler error", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}
	HandlerWithOptions(ctrl, options)

	return r
}
