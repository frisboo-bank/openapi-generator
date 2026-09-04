-- +goose Up
CREATE TABLE outbox (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL,
  subject text NOT NULL,
  data bytea NOT NULL,
  metadata bytea NOT NULL,
  queued_at timestamptz NOT NULL DEFAULT now(),
  published_at timestamptz NULL
);

-- +goose Down
DROP TABLE IF EXISTS outbox;
