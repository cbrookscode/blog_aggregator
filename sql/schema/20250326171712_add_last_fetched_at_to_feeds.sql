-- +goose Up
-- +goose StatementBegin
ALTER TABLE feeds
ADD last_fetched_at TIMESTAMP;
-- +goose StatementEnd
