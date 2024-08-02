package delete

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "rest_api_shortener/internal/lib/api/response"
	"rest_api_shortener/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("Alias is empty")

			render.JSON(w, r, resp.Error("Invalid request alias"))

			return
		}

		err := urlDeleter.DeleteURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("URL alias not found", "alias", alias)

				render.JSON(w, r, resp.Error("URL alias not found"))

				return
			}

			log.Info("Failed to delete URL alias", "alias", alias, "err", err)

			render.JSON(w, r, resp.Error("Failed to delete URL alias"))

			return
		}

		log.Info("URL alias deleted", "alias", alias)

		render.JSON(w, r, resp.OK())
	}
}
