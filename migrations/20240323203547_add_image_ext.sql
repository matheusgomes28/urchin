-- +goose Up
-- +goose StatementBegin
ALTER TABLE images ADD COLUMN ext TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE images DROP COLUMN ext;
-- +goose StatementEnd