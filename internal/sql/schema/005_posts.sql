-- +goose Up
CREATE TABLE posts (
	id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	title TEXT NOT NULL,
	url TEXT UNIQUE NOT NULL,
	description TEXT,
	published_at TIMESTAMP,
	feed_id BIGINT REFERENCES feed(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
