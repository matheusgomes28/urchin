# Urchin 🐚

Urchin is a headless CMS (Content Management System) written in Golang, designed to be fast, efficient, and easily extensible. It allows you to
create a website or blog, with any template you like, in only a few
commands.

![Really no head?](static/assets/nohead.gif "So no head meme?")

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

Once the database is up and migrated, you can run Urchin with your
favourite debugger setup. For the `Vscode` and `delve` setup, we
have provided the file `.vscode/launch.json` so you should just be
able to select the (admin) app from the `Vscode` debugging dropdown.

## Architecture

Currently, the architecture of `urchin` is still in its early days.
The plan is to have two main applications: the public facing application
to serve the content through a website, and the admin application that
can be hidden, where users can modify the settings, add posts, pages, etc.

![diagram of urchin's architecture](static/assets/urchin-architecture.png "Urchin Application Architecture")

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
admin_port = 8081 # The port in which the admin app will be running
image_dir = "./images" # Directory to use for storing uploaded images.

# Navbar section specifies the links that appear on
# the navbar
[navbar]
links = [
    { name = "Home", href = "/", title = "Homepage" },
    { name = "About", href = "/about", title = "About page" },
    { name = "Services", href = "/services", title = "Services page" },
    { name = "Images", href = "/images", title = "Images page" },
    { name = "Contact", href = "/contact", title = "Contacts page" },
]
```

**Important**: The configuration values above are used to start-up the local development database.

## Dependencies

For the development of Urchin, you require additional dependecies, that can easily be installed with go.

- [Templ](https://github.com/a-h/templ) (for generating Go files from temple-files)
- [Goose](https://github.com/pressly/goose) (for migrating the database that Urchin relies on)

*Optional*:

- [Air](https://github.com/cosmtrek/air) (for hot-reloading Go projects)

To install the development dependencies simply execute the following Go commands:

```bash
go install github.com/pressly/goose/v3/cmd/goose@v3.18.0 
go install github.com/a-h/templ/cmd/templ@v0.2.543 
go install github.com/cosmtrek/air@v1.49.0 
```

After installing the required dependecies and starting the pre-configured database, you can simply execute the following command to execute the migration of the database for development purposes.

```bash
source .dev.env # sets the environment variable for the goose command.
cd migrations/
GOOSE_DRIVER="mysql" GOOSE_DBSTRING="$MARIADB_USER:$MARIADB_PASSWORD@tcp($MARIADB_ADDRESS:$MARIADB_PORT)/$MARIADB_DATABASE" goose up
```
