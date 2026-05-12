--Generate: insert
INSERT INTO
	direct_channels (minor_id, major_id)
VALUES
	($1, $2)
RETURNING
	*;

--Generate: get_by_id
SELECT
	*
FROM
	direct_channels
WHERE
	minor_id = $1
	AND major_id = $2;

--Generate: delete
DELETE FROM direct_channels
WHERE
	minor_id = $1
	AND major_id = $2;
