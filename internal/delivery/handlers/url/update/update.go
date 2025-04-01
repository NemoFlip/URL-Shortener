package update

import (
	resp "RESTProject/internal/lib/api/response"
	"RESTProject/internal/lib/logger/sl"
	"RESTProject/internal/storage"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	Alias  string `json:"alias" validate:"required"`
	NewURL string `json:"new_url" validate:"required,url"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.3 --name=URLUpdater
type URLUpdater interface {
	UpdateURL(alias, newURL string) error
}

// @Summary      Update URL
// @Description  updates existing URL by alias
// @Tags         URL
// @Accept       json
// @Produce      json
// @Param        request body Request true "Request for updating url"
// @Success      200 {object} Response "Successfully updated url"
// @Failure      400 {object} Response "Invalid request"
// @Failure      404 {object} Response "URL not found"
// @Router       /url [put]
func New(log *slog.Logger, urlUpdater URLUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.update.New"
		log = log.With(
			slog.String("fn", fn),
		)

		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		if err := urlUpdater.UpdateURL(req.Alias, req.NewURL); err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("url not found", slog.String("alias", req.Alias))
				render.JSON(w, r, resp.Error("url not found"))
				return
			}

			log.Error("failed to update url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to update url"))
			return
		}

		log.Info("url updated successfully")

		responseOK(w, r, req.Alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
