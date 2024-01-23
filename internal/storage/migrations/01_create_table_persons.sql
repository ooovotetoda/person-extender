-- +goose Up
CREATE TABLE IF NOT EXISTS persons (
                                     id SERIAL PRIMARY KEY,
                                     name VARCHAR(100) NOT NULL,
                                     surname VARCHAR(100) NOT NULL,
                                     patronymic VARCHAR(100),
                                     age INT NOT NULL,
                                     gender VARCHAR(10) NOT NULL,
                                     country VARCHAR(5) NOT NULL
    );

-- +goose Down
DROP TABLE persons;