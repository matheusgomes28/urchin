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

## Running Urchin

Urchin is developed with Golang, so make sure you have a recent enough
version of the compiler and you also follow the instructions in the
following sections.

### Build Requirements

If you're runnig Urchin locally, you should install all the requirements needed
to build the application. Here's a list of all the dependencies needed:

- Goose for database migrations: `go install github.com/pressly/goose/v3/cmd/goose@v3.18.0`
- Templ for code generation: `go install github.com/a-h/templ/cmd/templ@v0.2.543`
- (optionally) Air for live reloading: `go install github.com/cosmtrek/air@v1.49.0`

Ensure that you have the binaries in the `$GOBIN` directory somewhere in your path,
so you can call these tools from the terminal.

Alternatively, if your platform supports `make`, run the following command from the
project repo:

```bash
make install-tools
```

### Database Migrations

Once the requirements are installed, make sure you run the Goose migrations for the database.
We recommend creating a database called `urchin` and running the following
command:

```bash
cd migrations
GOOSE_DRIVER="mysql" GOOSE_DBSTRING="root:root@/urchin" goose up
```

Replace the database connection string with the appropriate string
dependending on where your database is.

After you've replaced the default template files with your prefered
template, simply build and start the app with the following commands.

### Building and Running Urchin

If your platform has support for `Makefiles`, simply call `make`:

```bash
make build
./tmp/urchin --config urchin_config.toml
```

This will start Urchin on `http://localhost:8080`. You can change the
configuration by editing the `urchin_config.toml` file.

For more information, see the [configuration settings](#configuration).

## Running with Docker Compose

To run with `docker-compose`, use the following
command:

```bash
docker-compose -f docker/docker-compose.yml up
```

This will start two containers: one containing the `urchin` app,
serving on port `8080`, and another one serving the `mariadb`
database internally. This will also run the migrations automatically
to setup the database!

## Development

If you want to debug the application, you can use `docker compose`
to startup just the `mariadb` container, then hook Urchin to your
favourite debugger (e.g. Vscode).

To startup the `mariadb` database, run the following command from
the project root:

```sh
docker compose -f docker/mariadb.yml up
```

Wait a little bit for the database container to start, then run the
migration steps:

```sh
cd migrations
GOOSE_DRIVER="mysql" GOOSE_DBSTRING="root:root@/urchin" goose up
```

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

The runtime configuration can be done through a [toml](https://toml.io/en/) configuration file or by setting the mandatory environment variables (*fallback*). This approach was chosen because configuration via toml supports advanced features (i.e. *relationships*, *arrays*, etc.). The `.dev.env`-file is used to configure the development database through `docker-compose`.

### toml configuration

The application can be started by providing the `config` flag which has to be set to a toml configuration file. The file has to contain the following mandatory values:

```toml
database_address = "localhost" # Address to the MariaDB database
database_user = "urchin" # User to access database
database_password = "urchinpw" # Password for the database user
database_port = 3306 # The port to use for the application
database_name = "urchin" # The database to use for Urchin
webserver_port = 8080 # The application port Urchin should use
image_dir = "./images" # Directory to use for storing uploaded images.
```

**Important**: The configuration values above are used to start-up the local development database.

### Environment variables configuration (fallback)

If chosen, by setting the following environment variables the application can be started without providing a toml configuration file. 

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

To ease up the development process, Docker is highly recommended. This way you can use the `docker/mariadb.yml` to set up a predefined MariaDB database server. The docker-compose file references the `.dev.env` and creates the Urchin database and an application user.

```bash
$ docker-compose -f docker/mariadb.yml up -d
```

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
$ GOOSE_DRIVER="mysql" GOOSE_DBSTRING="$MARIADB_USER:$MARIADB_PASSWORD@tcp($MARIADB_ADDRESS:$MARIADB_PORT)/$MARIADB_DATABASE" goose up
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
      "args": [
          "--config",
          "urchin_config.toml"
      ]
    },
    {
      "name": "Urchin Admin",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/urchin-admin/main.go",
      "cwd": "${workspaceFolder}",
      "args": [
          "--config",
          "urchin_config.toml"
      ]
    }
  ]
}
```
However, the [Go-Extension](https://marketplace.visualstudio.com/items?itemName=golang.Go) must be installed before you can use these launch configurations.
