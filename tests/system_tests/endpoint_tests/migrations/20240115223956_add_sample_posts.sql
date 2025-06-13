-- +goose Up
-- +goose StatementBegin
INSERT INTO posts(title, content) VALUES(
    'My Very First Post',
    '# My Markdown File\n\n
## Subheading\n\n
This is a simple Markdown file with a heading and a subheading. Below is a code block example:\n\n
```python\n
print("Hello, Markdown!")\n
```');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM posts ORDER BY id DESC LIMIT 1;
-- +goose StatementEnd
