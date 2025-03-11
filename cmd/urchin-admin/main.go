package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	admin_app "github.com/matheusgomes28/urchin/admin-app"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/rs/zerolog/log"
)

func main() {
	// sets zerolog as the main logger
	// in this APP
	common.SetupLogger()

	config_toml := flag.String("config", "", "path to config toml file")
	flag.Parse()

	var app_settings common.AppSettings
	if *config_toml == "" {
		log.Error().Msgf("config not specified (--config)")
		os.Exit(-1)
	}

	log.Info().Msgf("reading config file %s", *config_toml)
	settings, err := common.ReadConfigToml(*config_toml)
	if err != nil {
		log.Error().Msgf("could not read toml: %v", err)
		os.Exit(-2)
	}
	app_settings = settings

	// Check for environment variable override for the HTTP port
	if portStr := os.Getenv("URCHIN_HTTP_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			log.Info().Msgf("using HTTP port from environment: %d", port)
			app_settings.AdminPort = port
		} else {
			log.Warn().Msgf("invalid URCHIN_HTTP_PORT value '%s', using config file value: %d", portStr, app_settings.AdminPort)
		}
	}

	database, err := database.MakeSqlConnection(
		app_settings.DatabaseUser,
		app_settings.DatabasePassword,
		app_settings.DatabaseAddress,
		app_settings.DatabasePort,
		app_settings.DatabaseName,
	)
	if err != nil {
		log.Fatal().Msgf("could not create database: %v", err)
		os.Exit(-1)
	}

	r := admin_app.SetupRoutes(app_settings, database)
	err = r.Run(fmt.Sprintf(":%d", app_settings.AdminPort))
	if err != nil {
		log.Error().Msgf("could not run app: %v", err)
		os.Exit(-1)
	}
}
