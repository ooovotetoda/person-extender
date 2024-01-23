package delete

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "person-extender/internal/lib/api/response"
	"person-extender/internal/lib/logger/sl"
	"strconv"
)

type PersonDeleter interface {
	DeletePerson(ID int64) error
}

func New(log *slog.Logger, personDeleter PersonDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.person.delete.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		ID := chi.URLParam(r, "id")
		if ID == "" {
			log.Error("ID is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		personID, err := strconv.ParseInt(ID, 10, 64)
		if err != nil {
			log.Error("failed to convert ID to int64", sl.Err(err))

			render.JSON(w, r, resp.Error("invalid ID format"))

			return
		}

		err = personDeleter.DeletePerson(personID)
		if err != nil {
			log.Error("failed to delete person", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("personnel deleted successfully")
	}
}
