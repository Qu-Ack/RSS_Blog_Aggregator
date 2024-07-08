-- +goose Up
CREATE TABLE feedfollow (
	id UUID PRIMARY KEY,
	feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	UNIQUE (user_id, feed_id)


	
);


-- +goose Down

DROP TABLE feedfollow;
