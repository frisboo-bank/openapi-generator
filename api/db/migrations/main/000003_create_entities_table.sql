-- +goose Up
CREATE TABLE IF NOT EXISTS entities (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  slug citext NOT NULL UNIQUE,
  name text NOT NULL,
  description text,
  version_lock int NOT NULL DEFAULT 1 CHECK (version_lock > 0),
  hidden_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz
);

CREATE INDEX idx_entity_slug ON entities (slug);

-- +goose Down
DROP TABLE IF EXISTS entities;

