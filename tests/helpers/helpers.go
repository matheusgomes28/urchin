package helpers

import (
	_ "database/sql"
	"embed"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"

	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
)

//go:generate ../../../migrations ./migrations

//go:embed migrations/*.sql
var EmbedMigrations embed.FS

func waitForDb(app_settings common.AppSettings) (database.SqlDatabase, error) {

	var err error = nil
	for range 400 {
		database, err := database.MakeMysqlConnection(app_settings.DatabaseUser, app_settings.DatabasePassword, app_settings.DatabaseAddress, app_settings.DatabasePort, app_settings.DatabaseName)

		if err == nil {
			return database, nil
		}

		time.Sleep(25 * time.Millisecond)
	}

	return database.SqlDatabase{}, fmt.Errorf("database did not start: %v", err)
}

func GetAppSettings() common.AppSettings {
	app_settings := common.AppSettings{
		DatabaseAddress:  "localhost",
		DatabasePort:     3306,
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseName:     "test",
		WebserverPort:    8080,
		CacheEnabled:     false,
	}

	return app_settings
}

func SetupDb(app_settings common.AppSettings) (cleanup func() error, db database.SqlDatabase, err error) {

	db, err = waitForDb(app_settings)
	if err != nil {
		return nil, database.SqlDatabase{}, fmt.Errorf("could not ping database: %v", err)
	}

	goose.SetBaseFS(EmbedMigrations)

	err = goose.SetDialect("mysql")
	if err != nil {
		return nil, database.SqlDatabase{}, fmt.Errorf("could not set the sql dialect: %v", err)
	}

	err = goose.Up(db.Connection, "migrations")
	if err != nil {
		return nil, database.SqlDatabase{}, fmt.Errorf("could not run migrations: %v", err)
	}

	cleanup = func() error {
		return goose.Reset(db.Connection, "migrations")
	}

	return cleanup, db, nil
}
