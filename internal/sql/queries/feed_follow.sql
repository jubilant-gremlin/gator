-- name: CreateFeedFollow :many
INSERT INTO feed_follows (created_at, updated_at, user_id, feed_id)
VALUES (
	$1,
	$2,
	$3,
	$4
	)
