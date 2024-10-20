package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	admin_app "github.com/matheusgomes28/urchin/admin-app"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/matheusgomes28/urchin/plugins"
	"github.com/rs/zerolog/log"
	lua "github.com/yuin/gopher-lua"
)

func loadShortcodeHandlers(shortcodes []common.Shortcode) (map[string]*lua.LState, error) {
	shortcode_handlers := make(map[string]*lua.LState, 0)
	for _, shortcode := range shortcodes {
		// Read the LUA state
		state := lua.NewState()
		err := state.DoFile(shortcode.Plugin)
		// TODO : check that the function HandleShortcode(args)
		//        exists and returns the correct type
		if err != nil {
			return map[string]*lua.LState{}, fmt.Errorf("could not load shortcode %s: %v", shortcode.Name, err)
		}
		shortcode_handlers[shortcode.Name] = state
	}

	return shortcode_handlers, nil
}

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

	shortcode_handlers, err := loadShortcodeHandlers(app_settings.Shortcodes)
	if err != nil {
		log.Error().Msgf("%v", err)
		os.Exit(-1)
	}

	// TODO : we probably want to refactor loadShortcodeHandler
	// TODO : into loadPluginHandlers instead

	post_hook := plugins.PostHook{}
	image_plugin := plugins.Plugin{
		ScriptName: "img",
		Id:         "img-plugin",
	}
	post_hook.Register(image_plugin)

	// img, _ := shortcode_handlers["img"]
	hooks_map := map[string]plugins.Hook{
		"add_post": post_hook,
	}

	r := admin_app.SetupRoutes(app_settings, database, shortcode_handlers, hooks_map)
	err = r.Run(fmt.Sprintf(":%d", app_settings.WebserverPort))
	if err != nil {
		log.Error().Msgf("could not run app: %v", err)
		os.Exit(-1)
	}
}
