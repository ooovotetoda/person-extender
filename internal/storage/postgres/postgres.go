package postgres

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"person-extender/internal/entity"
	"strings"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(host, port, user, password, dbname string) (*Storage, error) {
	const op = "storage.postgres.New"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
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

	storage := &Storage{db: db}

	err = goose.Up(storage.db, "internal/storage/migrations")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return storage, nil
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

func (s *Storage) GetPersons(filters *entity.Filters, limit, offset int64) ([]*entity.Person, error) {
	const op = "storage.postgres.GetPersons"

	query := "SELECT * FROM persons"

	conditions := []string{}
	params := []interface{}{}
	paramId := 1

	if filters.Name != nil {
		conditions = append(conditions, fmt.Sprintf("name = $%d", paramId))
		params = append(params, *filters.Name)
		paramId++
	}

	if filters.Surname != nil {
		conditions = append(conditions, fmt.Sprintf("surname = $%d", paramId))
		params = append(params, *filters.Surname)
		paramId++
	}

	if filters.Patronymic != nil {
		conditions = append(conditions, fmt.Sprintf("patronymic = $%d", paramId))
		params = append(params, *filters.Patronymic)
		paramId++
	}

	if filters.Age != nil {
		conditions = append(conditions, fmt.Sprintf("age = $%d", paramId))
		params = append(params, *filters.Age)
		paramId++
	}

	if filters.Gender != nil {
		conditions = append(conditions, fmt.Sprintf("gender = $%d", paramId))
		params = append(params, *filters.Gender)
		paramId++
	}

	if filters.Country != nil {
		conditions = append(conditions, fmt.Sprintf("country = $%d", paramId))
		params = append(params, *filters.Country)
		paramId++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := s.db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var persons []*entity.Person

	for rows.Next() {
		p := new(entity.Person)
		err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Country)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		persons = append(persons, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return persons, nil
}
