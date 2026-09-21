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

CREATE TRIGGER api_versions_updated_at
    BEFORE UPDATE ON apis
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS api_versions_updated_at ON api_versions;
DROP TABLE IF EXISTS api_versions;
