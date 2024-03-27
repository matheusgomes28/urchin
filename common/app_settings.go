package common

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type CardSchema struct {
	Name string `toml:"schema_name"`
}

type AppSettings struct {
	DatabaseAddress  string       `toml:"database_address"`
	DatabasePort     int          `toml:"database_port"`
	DatabaseUser     string       `toml:"database_user"`
	DatabasePassword string       `toml:"database_password"`
	DatabaseName     string       `toml:"database_name"`
	WebserverPort    int          `toml:"webserver_port"`
	ImageDirectory   string       `toml:"image_dir"`
	CardSchemas      []CardSchema `toml:"card_schema"`
}

func LoadSettings() (AppSettings, error) {
	// Want to load the environment variables
	database_address := os.Getenv("URCHIN_DATABASE_ADDRESS")
	if len(database_address) == 0 {
		return AppSettings{}, fmt.Errorf("URCHIN_DATABASE_ADDRESS is not defined")
	}

	database_user := os.Getenv("URCHIN_DATABASE_USER")
	if len(database_user) == 0 {
		return AppSettings{}, fmt.Errorf("URCHIN_DATABASE_USER is not defined")
	}

	database_password := os.Getenv("URCHIN_DATABASE_PASSWORD")
	if len(database_password) == 0 {
		return AppSettings{}, fmt.Errorf("URCHIN_DATABASE_PASSWORD is not defined")
	}

	database_name := os.Getenv("URCHIN_DATABASE_NAME")
	if len(database_name) == 0 {
		return AppSettings{}, fmt.Errorf("URCHIN_DATABASE_NAME is not defined")
	}

	database_port_str := os.Getenv("URCHIN_DATABASE_PORT")
	if len(database_port_str) == 0 {
		return AppSettings{}, fmt.Errorf("URCHIN_DATABASE_PORT is not defined")
	}

	database_port, err := strconv.Atoi(database_port_str)
	if err != nil {
		return AppSettings{}, fmt.Errorf("URCHIN_DATABASE_PORT is not a valid integer: %v", err)
	}

	webserver_port_str := os.Getenv("URCHIN_WEBSERVER_PORT")
	if webserver_port_str == "" {
		return AppSettings{}, fmt.Errorf("URCHIN_WEBSERVER_PORT is not defined")
	}

	webserver_port, err := strconv.Atoi(webserver_port_str)
	if err != nil {
		return AppSettings{}, fmt.Errorf("URCHIN_WEBSERVER_PORT is not valid: %v", err)
	}

	image_directory := os.Getenv("URCHIN_IMAGE_DIRECTORY")
	if len(image_directory) == 0 {
		return AppSettings{}, fmt.Errorf("URCHIN_IMAGE_DIRECTORY is not defined")
	}

	return AppSettings{
		DatabaseUser:     database_user,
		DatabasePassword: database_password,
		DatabaseAddress:  database_address,
		DatabasePort:     database_port,
		DatabaseName:     database_name,
		WebserverPort:    webserver_port,
		ImageDirectory:   image_directory,
	}, nil
}

func ReadConfigToml(filepath string) (AppSettings, error) {
	var config AppSettings
	metadata, err := toml.DecodeFile(filepath, &config)
	if err != nil {
		return AppSettings{}, err
	}

	if undecoded_keys := metadata.Undecoded(); len(undecoded_keys) > 0 {
		err := fmt.Errorf("could not decode keys: ")

		for _, key := range undecoded_keys {
			metadata.Keys()
			err = fmt.Errorf("%v %v,", err, strings.Join(key, "."))
		}

		return AppSettings{}, err
	}

	return config, nil
}
