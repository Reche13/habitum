-- +goose Up
-- +goose StatementBegin
ALTER TABLE habits 
ADD COLUMN current_streak INTEGER NOT NULL DEFAULT 0,
ADD COLUMN longest_streak INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE habits 
DROP COLUMN IF EXISTS current_streak,
DROP COLUMN IF EXISTS longest_streak;
-- +goose StatementEnd
