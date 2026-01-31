DO $$
DECLARE
    multiplier integer := COALESCE(NULLIF(current_setting('bloat.multiplier', true), ''), '1')::integer;
BEGIN
    FOR i IN 1..multiplier LOOP
        UPDATE bloatlab.orders
        SET order_status = 'paid'
        WHERE id % 3 = 0;

        UPDATE bloatlab.orders
        SET order_total = order_total + (random() * 100)::numeric(12, 2)
        WHERE id % 5 = 0;

        DELETE FROM bloatlab.orders
        WHERE id % 7 = 0;

        UPDATE bloatlab.events
        SET payload = md5(random()::text) || payload
        WHERE id % 4 = 0;

        DELETE FROM bloatlab.events
        WHERE id % 6 = 0;
    END LOOP;
END $$;

VACUUM (ANALYZE) bloatlab.orders;
VACUUM (ANALYZE) bloatlab.events;
