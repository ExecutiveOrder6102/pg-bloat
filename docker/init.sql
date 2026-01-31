CREATE SCHEMA IF NOT EXISTS bloatlab;

CREATE TABLE IF NOT EXISTS bloatlab.orders (
    id bigserial PRIMARY KEY,
    account_id integer NOT NULL,
    order_total numeric(12, 2) NOT NULL,
    order_status text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
) WITH (fillfactor = 70);

CREATE INDEX IF NOT EXISTS idx_orders_account_id ON bloatlab.orders (account_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON bloatlab.orders (order_status);

DO $$
DECLARE
    multiplier integer := COALESCE(NULLIF(current_setting('bloat.multiplier', true), ''), '1')::integer;
BEGIN
    INSERT INTO bloatlab.orders (account_id, order_total, order_status, created_at)
    SELECT
        (random() * 50000)::integer,
        (random() * 1000)::numeric(12, 2),
        (ARRAY['new', 'paid', 'shipped', 'cancelled'])[1 + (random() * 3)::integer],
        now() - (random() * 365)::integer * interval '1 day'
    FROM generate_series(1, 200000 * multiplier);
END $$;

CREATE TABLE IF NOT EXISTS bloatlab.events (
    id bigserial PRIMARY KEY,
    tenant_id integer NOT NULL,
    payload text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
) WITH (fillfactor = 70);

CREATE INDEX IF NOT EXISTS idx_events_tenant_id ON bloatlab.events (tenant_id);

DO $$
DECLARE
    multiplier integer := COALESCE(NULLIF(current_setting('bloat.multiplier', true), ''), '1')::integer;
BEGIN
    INSERT INTO bloatlab.events (tenant_id, payload, created_at)
    SELECT
        (random() * 10000)::integer,
        md5(random()::text),
        now() - (random() * 180)::integer * interval '1 day'
    FROM generate_series(1, 150000 * multiplier);
END $$;
