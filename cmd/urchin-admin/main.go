// @title        Urchin Admin API
// @version      1.0.0
// @description  This is the admin API for the Urchin app.
// @schemes   http
// @host      localhost:8081
// @BasePath  /
// @contact.name   MatheusGomes
// @contact.email  email@email.com
// @license.name  MIT
// @consumes  application/json
// @consumes  multipart/form-data
// @produces  application/json
package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	admin_app "github.com/matheusgomes28/urchin/admin-app"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"

	// Importa los documentos de Swagger generados para que estén disponibles.
	_ "github.com/matheusgomes28/urchin/docs/admin"
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

	database, err := database.MakeMysqlConnection(
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
	// Esta línea añade la ruta para la UI de Swagger.
	// La URL será: http://localhost:8081/swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err = r.Run(fmt.Sprintf(":%d", app_settings.AdminPort))
	if err != nil {
		log.Error().Msgf("could not run app: %v", err)
		os.Exit(-1)
	}
}
