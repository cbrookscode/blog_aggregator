// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: feeds.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const addFeed = `-- name: AddFeed :one
SELECT id, created_at, updated_at, name, url, last_fetched_at, user_id FROM feeds WHERE "name" = $1 and "url" = $2
`

type AddFeedParams struct {
	Name string
	Url  string
}

func (q *Queries) AddFeed(ctx context.Context, arg AddFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, addFeed, arg.Name, arg.Url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.LastFetchedAt,
		&i.UserID,
	)
	return i, err
}

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, "name", "url", user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING id, created_at, updated_at, name, url, last_fetched_at, user_id
`

type CreateFeedParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Url       string
	UserID    uuid.UUID
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Url,
		arg.UserID,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.LastFetchedAt,
		&i.UserID,
	)
	return i, err
}

const getFeedByUrl = `-- name: GetFeedByUrl :one
SELECT id, created_at, updated_at, name, url, last_fetched_at, user_id FROM feeds WHERE "url" = $1
`

func (q *Queries) GetFeedByUrl(ctx context.Context, url string) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeedByUrl, url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.LastFetchedAt,
		&i.UserID,
	)
	return i, err
}

const getFeeds = `-- name: GetFeeds :many
SELECT id, created_at, updated_at, name, url, last_fetched_at, user_id FROM feeds
`

func (q *Queries) GetFeeds(ctx context.Context) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.LastFetchedAt,
			&i.UserID,
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

const getNextFeedToFetch = `-- name: GetNextFeedToFetch :one
SELECT id, created_at, updated_at, name, url, last_fetched_at, user_id FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1
`

func (q *Queries) GetNextFeedToFetch(ctx context.Context) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getNextFeedToFetch)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.LastFetchedAt,
		&i.UserID,
	)
	return i, err
}

const markFeedFetched = `-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1
`

func (q *Queries) MarkFeedFetched(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, markFeedFetched, id)
	return err
}
