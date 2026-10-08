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

CREATE TABLE slug_history (
  slug citext NOT NULL,
  resource_type text NOT NULL CHECK (resource_type IN ('entity', 'api')),
  resource_id uuid,
  created_by citext NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (slug, resource_type)
);

CREATE INDEX idx_slug_history_lookup ON slug_history(slug);

-- +goose Down
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS slug_history;
