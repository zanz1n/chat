CREATE TABLE attachments (
	id bigserial PRIMARY KEY,
	created_at timestamp NOT NULL DEFAULT now(),
	format text NOT NULL,
	link text,
	data bytea,
	CONSTRAINT attachment_oneof_check CHECK (
		(link IS NOT NULL)::integer + (data IS NOT NULL)::integer = 1
	)
);

ALTER TABLE attachments
ALTER COLUMN data
SET
	storage external;

ALTER TABLE users
ADD CONSTRAINT users_attachment_fkey FOREIGN key (picture_id) REFERENCES attachments (id) ON UPDATE CASCADE ON DELETE SET NULL;

ALTER TABLE channels
ADD CONSTRAINT channels_attachment_fkey FOREIGN key (picture_id) REFERENCES attachments (id) ON UPDATE CASCADE ON DELETE SET NULL;

ALTER TABLE messages
ADD CONSTRAINT messages_attachment_fkey FOREIGN key (attachment_id) REFERENCES attachments (id) ON UPDATE CASCADE ON DELETE SET NULL;
