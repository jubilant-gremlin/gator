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
