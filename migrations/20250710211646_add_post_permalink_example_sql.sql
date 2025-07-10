-- +goose Up
-- +goose StatementBegin
INSERT INTO post_permalinks(permalink, post_id) VALUES("no_head", 2)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM post_permalinks;
-- +goose StatementEnd
