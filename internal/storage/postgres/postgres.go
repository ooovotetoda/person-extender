package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"person-extender/internal/entity"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

const (
	host = "localhost"
	port = 5432
)

func New() (*Storage, error) {
	const op = "storage.postgres.New"

	user, password, dbname :=
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SavePerson(person *entity.Person) (int64, error) {
	const op = "storage.postgres.SavePerson"

	stmt, err := s.db.Prepare("INSERT INTO persons (name, surname, patronymic, age, gender, country) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int64
	err = stmt.QueryRow(person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Country).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) DeletePerson(ID int64) error {
	const op = "storage.postgres.DeletePerson"

	_, err := s.db.Exec(`DELETE FROM persons WHERE id = $1`, ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdatePerson(person *entity.Person) error {
	const op = "storage.postgres.UpdatePerson"

	stmt, err := s.db.Prepare("UPDATE persons SET name = $2, surname = $3, patronymic = $4, age = $5, gender = $6, country = $7 WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(person.ID, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Country)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
