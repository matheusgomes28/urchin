package main

import (
	_ "github.com/go-sql-driver/mysql"
	admin_app "github.com/matheusgomes28/admin-app"
	"github.com/matheusgomes28/common"
	"github.com/matheusgomes28/database"
	"github.com/rs/zerolog/log"
)


func main() {

	// sets zerolog as the main logger
	// in this APP
	common.SetupLogger()

	app_settings, err := common.LoadSettings()
	if err != nil {
		log.Fatal().Msgf("could not load app settings: %v", err)
	}

	database, err := database.MakeSqlConnection(
		app_settings.Database_user,
		app_settings.Database_password,
		app_settings.Database_address,
		app_settings.Database_port,
		app_settings.Database_name,
	)
	if err != nil {
		log.Fatal().Msgf("could not create database: %v", err)
	}

	err = admin_app.Run(app_settings, database)
	if err != nil {
		log.Fatal().Msgf("could not run app: %v", err)
	}
}