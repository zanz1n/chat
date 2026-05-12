CREATE TYPE user_role AS enum('NONE', 'DEFAULT', 'MODERATOR', 'ADMIN');

CREATE TABLE users (
	id uuid PRIMARY KEY DEFAULT uuidv7(),
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now(),
	username text NOT NULL,
	name text,
	email text,
	pronouns text,
	picture_id bigint,
	bio text,
	role user_role NOT NULL,
	password bytea NOT NULL
);

CREATE UNIQUE INDEX users_username_idx ON users (username);

CREATE INDEX users_picture_idx ON users (picture_id);
