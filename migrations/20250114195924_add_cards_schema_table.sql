-- +goose Up
-- +goose StatementBegin
CREATE TABLE card_schemas (
  uuid BINARY(16) PRIMARY KEY NOT NULL,
  json_id VARCHAR(50) NOT NULL,
  json_schema JSON NOT NULL,
  json_title VARCHAR(50) NOT NULL,
  card_ids JSON
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE card_schemas;
-- +goose StatementEnd
