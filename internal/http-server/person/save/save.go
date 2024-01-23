package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	"person-extender/internal/entity"
	"person-extender/internal/lib/api"
	resp "person-extender/internal/lib/api/response"
	"person-extender/internal/lib/logger/sl"
)

type Request struct {
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Patronymic string `json:"patronymic,omitempty"`
}

type Response struct {
	resp.Response
	ID int64 `json:"id"`
}

type PersonSaver interface {
	SavePerson(person *entity.Person) (int64, error)
}

func New(log *slog.Logger, personSaver PersonSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.person.save.New"

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

		personExtends, err := api.GetPersonExtends(req.Name)
		if err != nil {
			log.Error("failed to get persons extends", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		person := &entity.Person{
			Name:       req.Name,
			Surname:    req.Surname,
			Patronymic: req.Patronymic,
			Age:        personExtends.Age,
			Gender:     personExtends.Gender,
			Country:    personExtends.Country,
		}

		ID, err := personSaver.SavePerson(person)
		if err != nil {
			log.Error("failed to save person", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("person successfully added", slog.Int64("id", ID))

		responseOK(w, r, ID)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, ID int64) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		ID:       ID,
	})
}
