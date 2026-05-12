--Generate: insert
INSERT INTO
	members (
		channel_id,
		user_id,
		added_at,
		updated_at,
		nickname,
		role
	)
VALUES
	($1, $2, $3, $4, $5, $6)
RETURNING
	*;

--Generate: get_by_id
SELECT
	*
FROM
	members
WHERE
	channel_id = $1
	AND user_id = $2;

--Generate: get_with_user
SELECT
	user.*,
	member.*
FROM
	members member
	JOIN users user ON member.user_id = user.id
WHERE
	member.channel_id = $1
	AND member.user_id = $2;

--Generate: get_by_channel
SELECT
	*
FROM
	members
WHERE
	channel_id = $1
	AND user_id < $2
LIMIT
	$3
ORDER BY
	user_id DESC;

--Generate: get_users_by_channel
SELECT
	user.*,
	member.*
FROM
	members member
	JOIN users user ON member.user_id = user.id
WHERE
	member.channel_id = $1
	AND member.user_id < $2
LIMIT
	$3
ORDER BY
	member.user_id DESC;

--Generate: update_nickname
UPDATE members
SET
	nickname = $3,
	updated_at = now()
WHERE
	channel_id = $1
	AND user_id = $2
RETURNING
	*;

--Generate: update_role
UPDATE members
SET role = $3,
updated_at = now()
WHERE
	channel_id = $1
	AND user_id = $2
RETURNING
	*;

--Generate: delete
DELETE FROM members
WHERE
	channel_id = $1
	AND user_id = $2
RETURNING
	*;
