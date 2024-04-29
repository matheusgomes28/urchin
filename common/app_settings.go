package common

import (
	"github.com/BurntSushi/toml"
)

type Navbar struct {
	Links []Link `toml:"links"`
}

type Shortcode struct {
	// name for the shortcode {{name:...:...}}
	Name string `toml:"name"`
	// The lua plugin path
	Plugin string `toml:"plugin"`
}

type AppSettings struct {
	DatabaseAddress  string      `toml:"database_address"`
	DatabasePort     int         `toml:"database_port"`
	DatabaseUser     string      `toml:"database_user"`
	DatabasePassword string      `toml:"database_password"`
	DatabaseName     string      `toml:"database_name"`
	WebserverPort    int         `toml:"webserver_port"`
	AdminPort        int         `toml:"admin_port"`
	ImageDirectory   string      `toml:"image_dir"`
	CacheEnabled     bool        `toml:"cache_enabled"`
	RecaptchaSiteKey string      `toml:"recaptcha_sitekey,omitempty"`
	RecaptchaSecret  string      `toml:"recaptcha_secret,omitempty"`
	AppNavbar        Navbar      `toml:"navbar"`
	Shortcodes       []Shortcode `toml:"shortcodes"`
}

func ReadConfigToml(filepath string) (AppSettings, error) {
	var config AppSettings
	_, err := toml.DecodeFile(filepath, &config)
	if err != nil {
		return AppSettings{}, err
	}

	return config, nil
}
