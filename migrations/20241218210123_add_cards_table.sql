-- +goose Up
-- +goose StatementBegin
CREATE TABLE cards (
  uuid VARCHAR(36) PRIMARY KEY DEFAULT UUID(),
  image_location TEXT NOT NULL,
  json_data TEXT NOT NULL,
  json_schema TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cards;
-- +goose StatementEnd
