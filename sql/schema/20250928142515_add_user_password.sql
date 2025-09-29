-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN password_hash TEXT NOT NULL DEFAULT 'unset';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN password_hash;
-- +goose StatementEnd
