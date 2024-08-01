package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "rest_api_shortener/internal/lib/api/responce"
	"rest_api_shortener/internal/storage"
)

// URLGetter is an interface for getting url by alias.
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.new"

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

		url, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("URL not found", "alias", alias)

				render.JSON(w, r, resp.Error("URL not found"))

				return
			}

			log.Error("Failed to get URL", "alias", alias, "err", err)

			render.JSON(w, r, resp.Error("Internal Error"))

			return
		}

		log.Info("Got URL", "alias", alias, "url", url)

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
