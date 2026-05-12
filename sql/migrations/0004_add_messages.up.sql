CREATE TABLE messages (
	id uuid PRIMARY KEY,
	-- mutually exclusive
	channel_id uuid,
	direct_id bigint,
	--
	user_id uuid,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	content text,
	attachment_id bigint,
	CONSTRAINT message_one_channel_check CHECK (
		(channel_id IS NOT NULL)::integer + (direct_id IS NOT NULL)::integer = 1
	),
	CONSTRAINT message_not_null_check CHECK (
		(content IS NOT NULL)
		OR (attachment_id IS NOT NULL)
	),
	CONSTRAINT messages_channel_fkey FOREIGN key (channel_id) REFERENCES channels (id) ON UPDATE CASCADE ON DELETE CASCADE,
	CONSTRAINT messages_direct_fkey FOREIGN key (direct_id) REFERENCES direct_channels (id) ON UPDATE CASCADE ON DELETE CASCADE,
	CONSTRAINT messages_user_fkey FOREIGN key (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE INDEX messages_channel_idx ON messages (channel_id);

CREATE INDEX messages_direct_idx ON messages (direct_id);

CREATE INDEX messages_user_idx ON messages (user_id);

CREATE INDEX messages_attachment_idx ON messages (attachment_id);

CREATE OR REPLACE FUNCTION message_notify () returns trigger AS $$
DECLARE
    payload_id text;
    payload json;
BEGIN
    IF (TG_OP = 'DELETE') THEN
        IF (OLD.channel_id IS NOT NULL) THEN
            payload_id = 'channels/' || cast(OLD.channel_id AS text);
        ELSE
            payload_id = 'direct/' || cast(OLD.direct_id AS text);
        END IF;

        payload = json_build_object(
            'action', TG_OP,
            'data', row_to_json(OLD)
        );
    ELSE
        IF (NEW.channel_id IS NOT NULL) THEN
            payload_id = 'channels/' || cast(NEW.channel_id AS text);
        ELSE
            payload_id = 'direct/' || cast(NEW.direct_id AS text);
        END IF;

        payload = json_build_object(
            'action', TG_OP,
            'data', row_to_json(NEW)
        );
    END IF;

    PERFORM pg_notify(payload_id, payload::text);

    RETURN NULL;
END
$$ language plpgsql;

CREATE TRIGGER message_create_trigger before insert ON messages FOR each ROW
EXECUTE procedure message_notify ();

CREATE TRIGGER message_update_trigger before
UPDATE ON messages FOR each ROW
EXECUTE procedure message_notify ();

CREATE TRIGGER message_delete_trigger before delete ON messages FOR each ROW
EXECUTE procedure message_notify ();
