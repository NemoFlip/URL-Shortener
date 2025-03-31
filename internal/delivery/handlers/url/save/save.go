package save

import (
	resp "RESTProject/internal/lib/api/response"
	"RESTProject/internal/lib/logger/sl"
	"RESTProject/internal/lib/random"
	"RESTProject/internal/storage"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

// TODO: move to the config
const aliasLength =  6

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.3 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave, alias string) error
}

// @Summary      Post URL
// @Description  posts url with alias
// @Tags         URL
// @Accept json
// @Produce json
// @Param request body Request true "Request for posting url"
// @Success      200  {object} Response "Successfully submited url"
// @Failure 400 {object} Response "Invalid request"
// @Router       /url [post]
func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.save.New"
		log = log.With(
			slog.String("fn", fn),
			//slog.String("request_id", middleware.GetReqID(r.Context())),
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

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		if err := urlSaver.SaveURL(req.URL, alias); err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url already exists", slog.String("url", req.URL))

				render.JSON(w, r, resp.Error("url already exists"))
				return
			}

			log.Error("failed to save url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to save url"))
			return
		}

		log.Info("url saved successfully")

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias: alias,
	})
}