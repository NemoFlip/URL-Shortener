package redirect

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

//go:generate go run github.com/vektra/mockery/v2@v2.53.3 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

// @Summary      Get URL
// @Description  redirect to url by speciefied alias
// @Tags         URL
// @Produce json
// @Param alias path string true "alias for further redirect"
// @Success      302  {object} string "Successfully redirected"
// @Failure 400 {object} string "Invalid request"
// @Router       /{alias} [get]
func New(logger *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.redirect.New"
		logger = logger.With(
			slog.String("fn", fn),
		)
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			logger.Info("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		url, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				logger.Info("url is not found by alias", slog.String("alias", alias))
				render.JSON(w, r, resp.Error("url not found"))
				return
			}

			logger.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		logger.Info("got url", slog.String("url", url))
		// redirect to found url
		http.Redirect(w, r, url, http.StatusFound)
	}

}
