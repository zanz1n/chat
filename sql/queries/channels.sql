--Generate: insert
INSERT INTO
	channels (owner_id, name, description)
VALUES
	($1, $2, $3)
RETURNING
	*;

--Generate: get_by_id
SELECT
	*
FROM
	channels
WHERE
	id = $1;

--Generate: get_by_user
SELECT
	*
FROM
	channels
WHERE
	owner_id = $1
	AND id < $2
LIMIT
	$3
ORDER BY
	id DESC;

--Generate: search_by_user
SELECT
	1;

--Generate: update
UPDATE channels
SET
	name = $2,
	description = $3,
	updated_at = now()
WHERE
	id = $1
RETURNING
	*;

--Generate: update_picture
UPDATE channels
SET
	picture_id = $2,
	updated_at = now()
WHERE
	id = $1;

--Generate: delete
DELETE FROM channels
WHERE
	id = $1
RETURNING
	*;
