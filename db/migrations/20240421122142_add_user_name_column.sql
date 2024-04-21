-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS message
ADD COLUMN IF NOT EXISTS from_user_name TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS message
DROP COLUMN IF EXISTS from_user_name;
-- +goose StatementEnd
