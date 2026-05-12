ALTER TABLE messages
DROP CONSTRAINT messages_attachment_fkey;

ALTER TABLE channels
DROP CONSTRAINT channels_attachment_fkey;

ALTER TABLE messages
DROP CONSTRAINT messages_attachment_fkey;

DROP TABLE IF EXISTS attachments;
