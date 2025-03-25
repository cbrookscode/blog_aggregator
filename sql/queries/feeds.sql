-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, "name", "url", user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: AddFeed :one
SELECT * FROM feeds WHERE "name" = $1 and "url" = $2;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: DeleteAllFeeds :exec
TRUNCATE feeds;