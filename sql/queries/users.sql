--Generate: insert
INSERT INTO
	users (username, name, email, pronouns, role, password)
VALUES
	($1, $2, $3, $4, $5, $6)
RETURNING
	*;

--Generate: get_by_id
SELECT
	*
FROM
	users
WHERE
	id = $1;

--Generate: get_by_username
SELECT
	*
FROM
	users
WHERE
	username = $1;

--Generate: get_many
SELECT
	*
FROM
	users
WHERE
	id < $1
LIMIT
	$2
ORDER BY
	id DESC;

--Generate: update_username
UPDATE users
SET
	username = $2,
	updated_at = now()
WHERE
	id = $1
RETURNING
	*;

--Generate: update
UPDATE users
SET
	pronouns = $2,
	name = $3,
	email = $4,
	bio = $5,
	updated_at = now()
WHERE
	id = $1
RETURNING
	*;

--Generate: update_picture
UPDATE users
SET
	picture_id = $2,
	updated_at = now()
WHERE
	id = $1;

--Generate: update_password
UPDATE users
SET
	password = $2,
	updated_at = now()
WHERE
	id = $1
RETURNING
	*;

--Generate: delete
DELETE FROM users
WHERE
	id = $1
RETURNING
	*;
