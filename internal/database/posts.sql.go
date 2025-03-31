// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const createPost = `-- name: CreatePost :exec
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
ON CONFLICT (url) DO NOTHING
`

type CreatePostParams struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description sql.NullString
	PublishedAt sql.NullTime
	FeedID      sql.NullInt64
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) error {
	_, err := q.db.ExecContext(ctx, createPost,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.FeedID,
	)
	return err
}

const getPostsForUser = `-- name: GetPostsForUser :many
SELECT id, created_at, updated_at, title, url, description, published_at, feed_id FROM posts
ORDER BY published_at DESC
LIMIT $1
`

func (q *Queries) GetPostsForUser(ctx context.Context, limit int32) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForUser, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const resetPosts = `-- name: ResetPosts :exec
DELETE FROM POSTS *
`

func (q *Queries) ResetPosts(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, resetPosts)
	return err
}
