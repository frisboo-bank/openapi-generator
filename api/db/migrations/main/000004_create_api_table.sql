-- +goose Up
CREATE TABLE IF NOT EXISTS apis (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  entity_id uuid NOT NULL REFERENCES entities(id),
  slug citext NOT NULL UNIQUE,
  name text NOT NULL,
  description text,
  version_lock int NOT NULL DEFAULT 1 CHECK (version_lock > 0),
  hidden_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz
);

CREATE INDEX idx_api_entity ON apis (entity_id);
CREATE INDEX idx_api_slug ON apis (slug);

-- +goose Down
DROP TABLE IF EXISTS apis;
