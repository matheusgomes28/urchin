package helpers

import (
	_ "database/sql"
	"embed"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
)

//go:generate ../../../migrations ./migrations

//go:embed migrations/*.sql
var EmbedMigrations embed.FS

func WaitForDb(app_settings common.AppSettings) (database.SqlDatabase, error) {

	for range 400 {
		database, err := database.MakeMysqlConnection(app_settings.DatabaseUser, app_settings.DatabasePassword, app_settings.DatabaseAddress, app_settings.DatabasePort, app_settings.DatabaseName)

		if err == nil {
			return database, nil
		}

		time.Sleep(25 * time.Millisecond)
	}

	return database.SqlDatabase{}, fmt.Errorf("database did not start")
}

// GetAppSettings gets the settings for the http servers
// taking into account a unique port. Very hacky way to
// get a unique port: manually have to pass a new number
// for every test...
// TODO : Find a way to assign a unique port at compile
//
//	time
func GetAppSettings() common.AppSettings {
	app_settings := common.AppSettings{
		DatabaseAddress:  "localhost",
		DatabasePort:     3306,
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseName:     "urchin",
		WebserverPort:    8080,
		CacheEnabled:     false,
	}

	return app_settings
}
