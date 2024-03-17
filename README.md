# Urchin 🐚

Urchin is a headless CMS (Content Management System) written in Golang, designed to be fast, efficient, and easily extensible. It allows you to
create a website or blog, with any template you like, in only a few
commands.

![Really no head?](static/nohead.gif "So no head meme?")

## Features 🚀

- [x] **Headless Architecture:** Adding pages, posts, or forms should all
  be done with easy requests to the API.
- [x] **Golang-Powered:** Leverage the performance and safety of one of the
  best languages in the market for backend development.
- [x] **SQL Database Integration:** Store your posts and pages in SQL databases for reliable and scalable data storage.
- [ ] **Post**: We can add, update, and delete posts. Posts can be served
  through a unique URL.
- [ ] **Pages**: TODO.
- [ ] **Menus**: TODO
- [ ] **Live Reload** through the use of `air`.

## Installation

Ensure you have Golang installed on your system before proceeding with the installation.

```bash
go get -u github.com/username/urchin
```

## Example - Running the App

First, ensure you have the neccesary libraries to run the application
```bash
make install-tools
```

Following that, make sure you run the Goose migrations for the database.
We recommend creating a database called `urchin` and running the following
command:

```bash
GOOSE_DRIVER="mysql" GOOSE_DBSTRING="root:root@/gocms" goose up
```

Replace the database connection string with the appropriate string
dependending on where your database is.

After you've replaced the default template files with your prefered
template, simply build and start the app with the following commands.

```bash
go build
./urchin
```

This will start Urchin on `http://localhost:8080`. You can customize
the configuration by providing the necessary environment variables.

For more information, see the [configuration settings](#configuration).

## Example - Running with Docker Compose

To run with `docker-compose`, use the following
command:

```bash
docker-compose up
```

This will start two containers: one containing the `urchin` app,
serving on port `8080`, and another one serving the `mariadb`
database internally. This will also run the migrations automatically
to setup the database!

## Architecture

Currently, the architecture of `urchin` is still in its early days.
The plan is to have two main applications: the public facing application
to serve the content through a website, and the admin application that
can be hidden, where users can modify the settings, add posts, pages, etc.

![diagram of urchin's architecture](static/urchin-architecture.png "Urchin Application Architecture")

In the above image, you can see the two applications running alongside,
and they share a database connection where the data is actually stored.
The list below explains some of the data intended to be stored in the
database:

- **posts**: a table where each row is an individual post, containing
  the title, content, and any other relevant data.
- **pages**: a table where HTML can be stored to be served as individual
  pages on a website.
- **cards**: Still TODO. Need to decide how this will allow users to display
  menu-like pages with cards.

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
- `URCHIN_DATABASE_PASSWORD` needs to contain the database password for the given user.
- `URCHIN_DATABASE_NAME` sets the name of the database Urchin will use.
- `URCHIN_WEBSERVER_PORT` the port the application should run on.
- `URCHIN_IMAGE_DIRECTORY` the directory images should be stored to if uploaded to Urchin

## Development

To ease up the development process, Docker is highly recommended. This way you can use the `docker/mysqldb.yml` to set up a predefined MySQL database server. The docker-compose file references the [`.dev.env`](#env-file) and creates the Urchin database and an application user.

```bash
$ docker-compose -f docker/mysqldb.yml up -d
```

### .env-file

As described under [configuration section](#configuration), specific environment variables have to be defined for Urchin to run properly. For this an `.dev.env` file is pre-configured to set the required variables with some dummy values. This .env-file includes variables for the database and Urchin and can be used for the launch configuration in VSCode.

### Dependencies

For the development of Urchin, you require additional dependecies, that can easily be installed with go.

- [Templ](https://github.com/a-h/templ) (for generating Go files from temple-files)
- [Goose](https://github.com/pressly/goose) (for migrating the database that Urchin relies on)

*Optional*:

- [Air](https://github.com/cosmtrek/air) (for hot-reloading Go projects)

To install the development dependencies simply execute the following Go commands:
```bash
$ go install github.com/pressly/goose/v3/cmd/goose@v3.18.0 
$ go install github.com/a-h/templ/cmd/templ@v0.2.543 
$ go install github.com/cosmtrek/air@v1.49.0 
```

After installing the required dependecies and starting the pre-configured database, you can simply execute the following command to execute the migration of the database for development purposes.

```bash
$ source .dev.env # sets the environment variable for the goose command.
$ cd migrations/
$ GOOSE_DRIVER="mysql" GOOSE_DBSTRING="$URCHIN_DATABASE_USER:$URCHIN_DATABASE_PASSWORD@tcp($URCHIN_DATABASE_ADDRESS:$URCHIN_DATABASE_PORT)/$URCHIN_DATABASE_NAME" goose up
```

### Launching & Debugging

To debug the application or simply running it from within VSCode create a `launch.json` with the following configuration: 

```json
{
  ...,
  "configurations": [
    {
      "name": "Urchin",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/urchin/main.go",
      "cwd": "${workspaceFolder}",
      "envFile": "${workspaceFolder}/.dev.env",
    },
    {
      "name": "Urchin Admin",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/urchin-admin/main.go",
      "cwd": "${workspaceFolder}",
      "envFile": "${workspaceFolder}/.dev.env",
    }
  ]
}
```
However, the [Go-Extension](https://marketplace.visualstudio.com/items?itemName=golang.Go) must be installed before you can use these launch configurations.