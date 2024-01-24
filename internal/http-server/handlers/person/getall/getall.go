package getall

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	"person-extender/internal/entity"
	resp "person-extender/internal/lib/api/response"
	"person-extender/internal/lib/logger/sl"
	"strconv"
)

type Request struct {
	Name       *string `json:"name,omitempty"`
	Surname    *string `json:"surname,omitempty"`
	Patronymic *string `json:"patronymic,omitempty"`
	Age        *int64  `json:"age,omitempty"`
	Gender     *string `json:"gender,omitempty"`
	Country    *string `json:"country,omitempty"`
}

type Response struct {
	resp.Response
	Persons []*entity.Person `json:"persons"`
}

type PersonsGetter interface {
	GetPersons(filters *entity.Filters, limit, offset int64) ([]*entity.Person, error)
}

func New(log *slog.Logger, personsGetter PersonsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.person.getall.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
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

		limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
		if err != nil {
			log.Error("failed to convert limit value", sl.Err(err))

			render.JSON(w, r, resp.Error("invalid limit value"))

			return
		}

		offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
		if err != nil {
			log.Error("failed to convert offset value", sl.Err(err))

			render.JSON(w, r, resp.Error("invalid offset value"))

			return
		}

		filters := &entity.Filters{
			Name:       req.Name,
			Surname:    req.Surname,
			Patronymic: req.Patronymic,
			Age:        req.Age,
			Gender:     req.Gender,
			Country:    req.Country,
		}

		persons, err := personsGetter.GetPersons(filters, limit, offset)

		log.Info("person successfully got")

		responseOK(w, r, persons)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, persons []*entity.Person) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Persons:  persons,
	})
}
