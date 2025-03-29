-- +goose Up
CREATE TABLE feed_follows (
	id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	feed_id BIGINT REFERENCES feed(id) ON DELETE CASCADE,
	UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;
