-- +goose Up
-- +goose StatementBegin
CREATE TABLE cards (
  uuid VARCHAR(36) DEFAULT(UUID()) PRIMARY KEY,
  image_location TEXT NOT NULL,
  json_data TEXT NOT NULL,
  json_schema TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cards;
-- +goose StatementEnd
