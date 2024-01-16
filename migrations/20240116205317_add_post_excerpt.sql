-- +goose Up
-- +goose StatementBegin
ALTER TABLE posts ADD COLUMN excerpt TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE posts DROP COLUMN excerpt;
-- +goose StatementEnd
