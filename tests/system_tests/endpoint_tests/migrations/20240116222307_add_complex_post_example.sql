-- +goose Up
-- +goose StatementBegin
INSERT INTO posts(title, content, excerpt) VALUES(
    'My Very First Post',
    '# Urchin 🐚

Urchin is a headless CMS (Content Management System) written in Golang, designed to be fast, efficient, and easily extensible. It allows you to
create a website or blog, with any template you like, in only a few
commands.

## Features 🚀

- **Headless Architecture:** Adding pages, posts, or forms should all
  be done with easy requests to the API.
- **Golang-Powered:** Leverage the performance and safety of one of the
  best languages in the market for backend development.
- **SQL Database Integration:** Store your posts and pages in SQL databases for reliable and scalable data storage.

## Installation

Ensure you have Golang installed on your system before proceeding with the installation.

```bash
go get -u github.com/username/urchin
```

## Example - Running the App

After you\'ve replaced the default template files with your prefered
template, simply build and start the app with the following commands.

```bash
go build
./urchin
```

This will start Urchin on `http://localhost:8080`. You can customize
the configuration by providing the necessary environment variables.

For more information, see the [configuration settings](#configuration).

## Dependencies

Urchin relies on the following Golang dependencies:

- [Gin](github.com/gin-gonic/gin) as the web framework for Golang.
- [ZeroLog](https://github.com/rs/zerolog) for logging.

## Configuration

The runtime configuration is handled through reading the
necessary environment variables. This approach was chosen as
it makes integrating `envfile`s quite easy.

The following list outlines the environment variables needed.

- `URCHIN_DATABASE_ADDRESS` should contain the database addres,
  e.g. `localhost`.
- `URCHIN_DATABASE_PORT` should be the connection port to the
  database. For example `3306`.
- `URCHIN_DATABASE_USER` is the database username.
- `URCHIN_DATABASE_PASSWORD` needs to contain the database
  password for the given user.

## License

Urchin is released under the MIT License. See LICENSE (TODO) for
details. Feel free to fork, modify, and use it in your projects!',
'This is Urchin! This post is an example of how markdown can be rendered as a post.');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM posts ORDER BY id DESC LIMIT 1;
-- +goose StatementEnd
