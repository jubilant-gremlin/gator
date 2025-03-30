-- name: CreateFeedEntry :one

INSERT INTO feed(name, url, created_at, updated_at, user_id)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
)
RETURNING *;

-- name: ResetFeed :exec

DELETE FROM feed *;

-- name: GetFeeds :many

SELECT name, url, user_id FROM feed;

-- name: GetFeed :one

SELECT * FROM feed WHERE url = $1;

-- name: MarkFeedFetched :exec

UPDATE feed SET last_fetched_at = $1, updated_at = $1
WHERE feed.id = $2;

-- name: GetNextFeedToFetch :one
SELECT * FROM feed
ORDER BY last_fetched_at ASC NULLS FIRST;
