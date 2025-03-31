package delete

import (
	resp "RESTProject/internal/lib/api/response"
	"RESTProject/internal/lib/logger/sl"
	"RESTProject/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

// @Summary      Get URL
// @Description  redirect to url by speciefied alias
// @Tags         URL
// @Produce json
// @Param alias path string true "alias for further redirect"
// @Success      200  {object} string "Successfully deleted"
// @Failure 400 {object} string "Invalid request"
// @Router       /url/{alias} [delete]
func New(logger *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.delete.New"
		logger = logger.With(
			slog.String("fn", fn),
		)
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			logger.Info("alias is absent")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		if err := urlDeleter.DeleteURL(alias); err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				logger.Info("url is not found by alias", slog.String("alias", alias))
				render.JSON(w, r, resp.Error("url not found"))
				return
			}

			logger.Error("failed to delete url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to delete url"))
			return
		}
		logger.Info("successfully deleted url")

		render.JSON(w, r, resp.OK())
	}
}
