CREATE TABLE channels (
	id uuid PRIMARY KEY,
	owner_id uuid NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now(),
	name text NOT NULL,
	description text,
	picture_id bigint,
	CONSTRAINT channels_owner_fkey FOREIGN key (owner_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX channels_owner_idx ON channels (owner_id);

CREATE INDEX channels_picture_idx ON channels (picture_id);

CREATE TABLE members (
	channel_id uuid NOT NULL,
	user_id uuid NOT NULL,
	added_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	nickname text NOT NULL,
	role integer NOT NULL,
	CONSTRAINT members_channel_fkey FOREIGN key (channel_id) REFERENCES channels (id) ON UPDATE CASCADE ON DELETE CASCADE,
	CONSTRAINT members_user_fkey FOREIGN key (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
	PRIMARY KEY (channel_id, user_id)
);

CREATE INDEX members_channel_idx ON members (channel_id);

CREATE INDEX members_user_idx ON members (user_id);

CREATE TABLE direct_channels (
	id bigserial PRIMARY KEY,
	minor_id uuid NOT NULL,
	major_id uuid NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now(),
	CONSTRAINT direct_minor_fkey FOREIGN key (minor_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
	CONSTRAINT direct_major_fkey FOREIGN key (major_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE UNIQUE INDEX direct_users_idx ON direct_channels (minor_id, major_id);
