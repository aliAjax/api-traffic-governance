CREATE TABLE IF NOT EXISTS traffic_services (id text PRIMARY KEY, tenant_id text NOT NULL, name text NOT NULL, region text, created_at timestamptz NOT NULL);
CREATE TABLE IF NOT EXISTS traffic_routes (id text PRIMARY KEY, service_id text NOT NULL, version integer NOT NULL, state text NOT NULL, created_at timestamptz NOT NULL);
CREATE TABLE IF NOT EXISTS traffic_quotas (id text PRIMARY KEY, tenant_id text NOT NULL, quota_key text NOT NULL, rate_per_second integer NOT NULL, burst integer NOT NULL, updated_at timestamptz NOT NULL);
CREATE UNIQUE INDEX IF NOT EXISTS traffic_quota_tenant_key ON traffic_quotas(tenant_id,quota_key);
