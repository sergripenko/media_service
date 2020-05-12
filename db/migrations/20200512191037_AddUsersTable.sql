
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email       VARCHAR(30) NOT NULL,
    first_name  VARCHAR(30),
    last_name   VARCHAR(30),
    created_at  TIMESTAMP,
    updated_at  TIMESTAMP,
    deleted_at  TIMESTAMP
);

INSERT INTO users(email, first_name, last_name) VALUES ('sergripenko@gmail.com', 'Sergey', 'Ripenko');

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS users;
