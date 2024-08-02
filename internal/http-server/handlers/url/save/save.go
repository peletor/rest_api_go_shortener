package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	resp "rest_api_shortener/internal/lib/api/response"
	"rest_api_shortener/internal/lib/random"
	"rest_api_shortener/internal/storage"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	Response resp.Response
	Alias    string `json:"alias,omitempty"`
}

// TODO: move to config
const randomAliasLength = 6

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Filed to decode request body", slog.Any("error", err))

			render.JSON(w, r, resp.Error("Filed to decode request"))

			return
		}

		log.Info("Request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErrors := err.(validator.ValidationErrors)

			log.Error("Invalid request", slog.Any("error", err))

			render.JSON(w, r, resp.ValidationError(validateErrors))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(randomAliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("URL already exists", slog.String("url", req.URL))

				render.JSON(w, r, resp.Error("URL already exists"))

				return
			}

			log.Error("Failed to save URL", slog.Any("error", err))

			render.JSON(w, r, resp.Error("Failed to save URL"))

			return
		}

		log.Info("URL saved",
			slog.String("url", req.URL),
			slog.Int64("id", id),
		)

		responceOK(w, r, alias)
	}
}

func responceOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})

}
