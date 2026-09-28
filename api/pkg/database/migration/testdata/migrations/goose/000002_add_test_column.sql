-- +goose Up
ALTER TABLE test_migrations ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- +goose Down
ALTER TABLE test_migrations DROP COLUMN IF EXISTS created_at;
