-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS film_comments (
    id SERIAL PRIMARY KEY,
    film_id INTEGER NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    comment TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_film_comments_film_id FOREIGN KEY (film_id) REFERENCES film(film_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS film_comments;
-- +goose StatementEnd
