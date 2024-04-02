package common

import (
	"github.com/BurntSushi/toml"
)

type AppSettings struct {
	DatabaseAddress  string `toml:"database_address"`
	DatabasePort     int    `toml:"database_port"`
	DatabaseUser     string `toml:"database_user"`
	DatabasePassword string `toml:"database_password"`
	DatabaseName     string `toml:"database_name"`
	WebserverPort    int    `toml:"webserver_port"`
	ImageDirectory   string `toml:"image_dir"`
}

func ReadConfigToml(filepath string) (AppSettings, error) {
	var config AppSettings
	_, err := toml.DecodeFile(filepath, &config)
	if err != nil {
		return AppSettings{}, err
	}

	return config, nil
}
