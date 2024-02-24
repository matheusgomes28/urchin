package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/matheusgomes28/urchin/app"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/rs/zerolog/log"
)

func main() {
	/// Load global application settings
	app_settings, err := common.LoadSettings()
	if err != nil {
		log.Error().Msgf("could not get app settings: %v\n", err)
		os.Exit(-1)
	}

	db_connection, err := database.MakeSqlConnection(
		app_settings.DatabaseUser,
		app_settings.DatabasePassword,
		app_settings.DatabaseAddress,
		app_settings.DatabasePort,
		app_settings.DatabaseName,
	)
	if err != nil {
		log.Error().Msgf("could not create database connection: %v", err)
		os.Exit(-1)
	}

	r := app.SetupRoutes(app_settings, db_connection)
	err = r.Run(fmt.Sprintf(":%d", app_settings.WebserverPort))
	if err != nil {
		log.Error().Msgf("could not run app: %v", err)
		os.Exit(-1)
	}
}
