
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE images (
    id          SERIAL PRIMARY KEY,
    user_id     INT NOT NULL,
    filename    VARCHAR(50),
    height      INT,
    width       INT,
    uniq_id     VARCHAR(50),
    url         VARCHAR(255),
    is_orig_img BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMP,
    updated_at  TIMESTAMP,
    deleted_at  TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS images;
