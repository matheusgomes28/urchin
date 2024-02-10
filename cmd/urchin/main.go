package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/matheusgomes28/app"
	"github.com/matheusgomes28/common"
	"github.com/matheusgomes28/database"
	"github.com/rs/zerolog/log"
)

func main() {
	/// Load global application settings
	app_settings, err := common.LoadSettings()
	if err != nil {
		log.Error().Msgf("could not get app settings: %v\n", err)
		return
	}

	db_connection, err := database.MakeSqlConnection(
		app_settings.Database_user,
		app_settings.Database_password,
		app_settings.Database_address,
		app_settings.Database_port,
		app_settings.Database_name,
	)
	if err != nil {
		log.Error().Msgf("could not create database connection: %v", err)
	}

	app.Run(app_settings, &db_connection)
}
