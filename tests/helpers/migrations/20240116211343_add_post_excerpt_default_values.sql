-- +goose Up
-- +goose StatementBegin
UPDATE posts SET excerpt = 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Praesent sed auctor neque, in interdum nisi. Duis pulvinar risus eu placerat feugiat. Morbi blandit bibendum molestie.';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- 
UPDATE posts SET excerpt = NULL;
-- +goose StatementEnd
