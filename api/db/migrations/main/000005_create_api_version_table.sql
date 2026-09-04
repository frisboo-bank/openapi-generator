-- +goose Up
CREATE TABLE IF NOT EXISTS api_versions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  api_id uuid NOT NULL REFERENCES apis(id),
  version text NOT NULL,
  name text NOT NULL,
  description text,
  base_path text,
  version_lock int NOT NULL DEFAULT 1 CHECK (version_lock > 0),
  hidden_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz,
  UNIQUE(api_id, version)
);

CREATE INDEX idx_api_version_api ON api_versions (api_id);
CREATE INDEX idx_api_version_version ON api_versions (version);

-- +goose Down
DROP TABLE IF EXISTS api_versions;
