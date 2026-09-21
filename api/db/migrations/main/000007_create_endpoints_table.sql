-- +goose Up
CREATE TABLE IF NOT EXISTS endpoints(
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  api_version_id  uuid NOT NULL REFERENCES api_versions(id) ON DELETE CASCADE,
  operation_id    text,
  method          text NOT NULL CHECK (method IN ('GET','POST','PUT','PATCH','DELETE','OPTIONS','HEAD')),
  path            text NOT NULL,
  summary         text,
  description     text,
  deprecated      boolean NOT NULL DEFAULT false,
  version_lock    int NOT NULL DEFAULT 1 CHECK (version_lock > 0),
  hidden_at       timestamptz,
  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz,
  UNIQUE(api_version_id, method, path),
  UNIQUE(api_version_id, operation_id)
);

CREATE TABLE IF NOT EXISTS endpoint_headers(
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  endpoint_id uuid NOT NULL REFERENCES endpoints(id) ON DELETE CASCADE,
  name        text NOT NULL,
  description text,
  value       text NOT NULL,
  required    boolean NOT NULL DEFAULT false,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz,
  UNIQUE(endpoint_id, name)
);

CREATE TABLE IF NOT EXISTS endpoint_parameters(
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  endpoint_id uuid NOT NULL REFERENCES endpoints(id) ON DELETE CASCADE,
  name        text NOT NULL,
  description text,
  param_in    text NOT NULL CHECK (param_in IN ('path','query','header','cookie')),
  required    boolean NOT NULL DEFAULT false,
  param_type  text NOT NULL DEFAULT 'string',
  format      text,
  default_val jsonb,
  example     jsonb,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz,
  UNIQUE(endpoint_id, name, param_in)
);

CREATE TABLE IF NOT EXISTS endpoint_request_bodies(
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  endpoint_id  uuid NOT NULL REFERENCES endpoints(id) ON DELETE CASCADE,
  description  text,
  required     boolean NOT NULL DEFAULT false,
  content_type text NOT NULL DEFAULT 'application/json',
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz,
  UNIQUE(endpoint_id, content_type)
);

CREATE TABLE IF NOT EXISTS endpoint_body_properties (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  body_id     uuid NOT NULL REFERENCES endpoint_request_bodies(id) ON DELETE CASCADE,
  name        text NOT NULL,
  description text,
  prop_type   text NOT NULL DEFAULT 'string',
  format      text,
  required    boolean NOT NULL DEFAULT false,
  example     jsonb,
  parent_id   uuid REFERENCES endpoint_body_properties(id) ON DELETE CASCADE,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz,
  UNIQUE(body_id, name, parent_id)
);

CREATE TABLE IF NOT EXISTS endpoint_responses (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  endpoint_id  uuid NOT NULL REFERENCES endpoints(id) ON DELETE CASCADE,
  status_code  text NOT NULL,
  description  text,
  content_type text NOT NULL DEFAULT 'application/json',
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz,
  UNIQUE(endpoint_id, status_code, content_type)
);

CREATE TABLE IF NOT EXISTS endpoint_response_properties (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  response_id uuid NOT NULL REFERENCES endpoint_responses(id) ON DELETE CASCADE,
  name        text NOT NULL,
  prop_type   text NOT NULL DEFAULT 'string',
  format      text,
  description text,
  required    boolean NOT NULL DEFAULT false,
  example     jsonb,
  parent_id   uuid REFERENCES endpoint_response_properties(id) ON DELETE CASCADE,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz,
  UNIQUE(response_id, name, parent_id)
);

CREATE TABLE IF NOT EXISTS endpoint_security (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  endpoint_id uuid NOT NULL REFERENCES endpoints(id) ON DELETE CASCADE,
  scheme_name text NOT NULL,
  scopes      text[],
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz,
  UNIQUE(endpoint_id, scheme_name)
);

CREATE TABLE IF NOT EXISTS endpoint_tags (
  endpoint_id uuid NOT NULL REFERENCES endpoints(id) ON DELETE CASCADE,
  tag         citext NOT NULL,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz,
  PRIMARY KEY (endpoint_id, tag)
);

CREATE INDEX idx_endpoints_api_version ON endpoints(api_version_id);
CREATE INDEX idx_endpoints_operation_id ON endpoints(api_version_id, operation_id) WHERE operation_id IS NOT NULL;
CREATE INDEX idx_endpoint_headers_endpoint ON endpoint_headers(endpoint_id);
CREATE INDEX idx_endpoint_parameters_endpoint ON endpoint_parameters(endpoint_id);
CREATE INDEX idx_endpoint_body_props_body ON endpoint_body_properties(body_id);
CREATE INDEX idx_endpoint_body_props_parent ON endpoint_body_properties(parent_id);
CREATE INDEX idx_endpoint_resp_props_response ON endpoint_response_properties(response_id);
CREATE INDEX idx_endpoint_resp_props_parent ON endpoint_response_properties(parent_id);

CREATE TRIGGER endpoints_updated_at BEFORE UPDATE ON endpoints FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER endpoint_headers_updated_at BEFORE UPDATE ON endpoint_headers FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER endpoint_parameters_updated_at BEFORE UPDATE ON endpoint_parameters FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER endpoint_request_bodies_updated_at BEFORE UPDATE ON endpoint_request_bodies FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER endpoint_body_properties_updated_at BEFORE UPDATE ON endpoint_body_properties FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER endpoint_responses_updated_at BEFORE UPDATE ON endpoint_responses FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER endpoint_response_properties_updated_at BEFORE UPDATE ON endpoint_response_properties FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER endpoint_security_updated_at BEFORE UPDATE ON endpoint_security FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER endpoint_tags_updated_at BEFORE UPDATE ON endpoint_tags FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS endpoint_tags_updated_at ON endpoint_tags;
DROP TRIGGER IF EXISTS endpoint_security_updated_at ON endpoint_security;
DROP TRIGGER IF EXISTS endpoint_response_properties_updated_at ON endpoint_response_properties;
DROP TRIGGER IF EXISTS endpoint_responses_updated_at ON endpoint_responses;
DROP TRIGGER IF EXISTS endpoint_body_properties_updated_at ON endpoint_body_properties;
DROP TRIGGER IF EXISTS endpoint_request_bodies_updated_at ON endpoint_request_bodies;
DROP TRIGGER IF EXISTS endpoint_parameters_updated_at ON endpoint_parameters;
DROP TRIGGER IF EXISTS endpoint_headers_updated_at ON endpoint_headers;
DROP TRIGGER IF EXISTS endpoints_updated_at ON endpoints;

DROP TABLE IF EXISTS endpoint_tags;
DROP TABLE IF EXISTS endpoint_security;
DROP TABLE IF EXISTS endpoint_response_properties;
DROP TABLE IF EXISTS endpoint_responses;
DROP TABLE IF EXISTS endpoint_body_properties;
DROP TABLE IF EXISTS endpoint_request_bodies;
DROP TABLE IF EXISTS endpoint_parameters;
DROP TABLE IF EXISTS endpoint_headers;
DROP TABLE IF EXISTS endpoints;
