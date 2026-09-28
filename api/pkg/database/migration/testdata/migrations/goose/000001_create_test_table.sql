-- +goose Up
CREATE TABLE test_migrations (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS test_migrations;
