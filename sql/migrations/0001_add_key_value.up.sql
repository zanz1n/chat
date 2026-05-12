CREATE EXTENSION if NOT EXISTS pg_cron;

CREATE UNLOGGED TABLE key_value (
	key text PRIMARY KEY,
	expiration timestamp,
	value jsonb
);

CREATE INDEX key_value_expiration_idx ON key_value (expiration);

CREATE OR REPLACE FUNCTION key_value_clear () returns void AS $$
BEGIN
    DELETE FROM key_value
    WHERE expiration IS NOT NULL AND expiration < now();
    -- VACUUM;
END;
$$ language plpgsql;

CREATE OR REPLACE FUNCTION key_value_enable () returns void AS $$
BEGIN
    PERFORM 1 FROM key_value WHERE key = '__key_value_enabled';
    IF NOT FOUND THEN
        SELECT cron.schedule(
            'cron_key_value_clear',
            '0 */1 * * *',
            'SELECT key_value_clear()'
        );
        INSERT INTO key_value(key) VALUES ('__key_value_enabled');
    END IF;
END;
$$ language plpgsql;

CREATE OR REPLACE FUNCTION key_value_disable () returns void AS $$
BEGIN
    PERFORM 1 FROM key_value WHERE key = '__key_value_enabled';
    IF FOUND THEN
        DELETE FROM key_value WHERE key = '__key_value_enabled';
        SELECT cron.unschedule('cron_key_value_clear');
        DELETE FROM key_value;
        -- VACUUM;
    END IF;
END;
$$ language plpgsql;
