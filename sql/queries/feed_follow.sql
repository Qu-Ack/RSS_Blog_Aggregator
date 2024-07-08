-- name: CreateFeedFollow :one
INSERT INTO feedfollow(id, feed_id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: DeleteFeedFollow :one
DELETE FROM feedfollow WHERE id=$1 RETURNING *;


-- name: GetAllFeedFollowOfUser :many
SELECT * FROM feedfollow WHERE user_id=$1;
