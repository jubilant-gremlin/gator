-- name: CreatePost :exec
INSERT INTO posts (created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7
	)
ON CONFLICT (url) DO NOTHING;

-- name: GetPostsForUser :many
SELECT * FROM posts
ORDER BY published_at DESC
LIMIT $1;

-- name: ResetPosts :exec
DELETE FROM POSTS *;
