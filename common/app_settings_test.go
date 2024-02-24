package common

import (
	"errors"
	"os"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
)

// Writes the contents into a temporary
// toml file
func writeToml(contents []byte) (s string, err error) {
	file, err := os.CreateTemp(os.TempDir(), "*.toml")
	if err != nil {
		return "", err
	}
	defer func() {
		err = errors.Join(file.Close(), err)
	}()

	_, err = file.Write(contents)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

func TestCorrectToml(t *testing.T) {
	expected := AppSettings{
		DatabaseAddress:  "test_database_address",
		DatabaseUser:     "test_database_user",
		DatabasePassword: "test_database_password",
		DatabaseName:     "test_database_name",
		WebserverPort:    99999,
		DatabasePort:     666,
	}
	bytes, err := toml.Marshal(expected)
	assert.Nil(t, err)

	filepath, err := writeToml(bytes)
	assert.Nil(t, err)

	actual, err := ReadConfigToml(filepath)
	assert.Nil(t, err)
	assert.Equal(t, actual, expected)
}

func TestMissingDatabaseAddress(t *testing.T) {

	missing_database_address := struct {
		DatabaseUser     string `toml:"database_user"`
		DatabasePassword string `toml:"database_password"`
		DatabaseName     string `toml:"database_name"`
		WebserverPort    string `toml:"webserver_port"`
		DatabasePort     string `toml:"database_port"`
	}{
		DatabaseUser:     "test_database_user",
		DatabasePassword: "test_database_password",
		DatabaseName:     "test_database_name",
		WebserverPort:    "99999",
		DatabasePort:     "666",
	}

	bytes, err := toml.Marshal(missing_database_address)
	assert.Nil(t, err)

	filepath, err := writeToml(bytes)
	assert.Nil(t, err)

	_, err = ReadConfigToml(filepath)
	assert.NotNil(t, err)
}

func TestMissingDatabaseUser(t *testing.T) {

	missing_database_user := struct {
		DatabaseAddress  string `toml:"database_address"`
		DatabasePassword string `toml:"database_password"`
		DatabaseName     string `toml:"database_name"`
		WebserverPort    string `toml:"webserver_port"`
		DatabasePort     string `toml:"database_port"`
	}{
		DatabaseAddress:  "test_database_address",
		DatabasePassword: "test_database_password",
		DatabaseName:     "test_database_name",
		WebserverPort:    "99999",
		DatabasePort:     "666",
	}

	bytes, err := toml.Marshal(missing_database_user)
	assert.Nil(t, err)

	filepath, err := writeToml(bytes)
	assert.Nil(t, err)

	_, err = ReadConfigToml(filepath)
	assert.NotNil(t, err)
}

func TestMissingDatabasePassword(t *testing.T) {
	missing_database_password := struct {
		DatabaseAddress string `toml:"database_address"`
		DatabaseUser    string `toml:"database_user"`
		DatabaseName    string `toml:"database_name"`
		WebserverPort   string `toml:"webserver_port"`
		DatabasePort    string `toml:"database_port"`
	}{
		DatabaseAddress: "test_database_address",
		DatabaseUser:    "test_database_user",
		DatabaseName:    "test_database_name",
		WebserverPort:   "99999",
		DatabasePort:    "666",
	}

	bytes, err := toml.Marshal(missing_database_password)
	assert.Nil(t, err)

	filepath, err := writeToml(bytes)
	assert.Nil(t, err)

	_, err = ReadConfigToml(filepath)
	assert.NotNil(t, err)
}

func TestMissingDatabaseName(t *testing.T) {
	missing_database_name := struct {
		DatabaseAddress  string `toml:"database_address"`
		DatabaseUser     string `toml:"database_user"`
		DatabasePassword string `toml:"database_password"`
		WebserverPort    string `toml:"webserver_port"`
		DatabasePort     string `toml:"database_port"`
	}{
		DatabaseAddress:  "test_database_address",
		DatabaseUser:     "test_database_user",
		DatabasePassword: "test_database_password",
		WebserverPort:    "99999",
		DatabasePort:     "666",
	}

	bytes, err := toml.Marshal(missing_database_name)
	assert.Nil(t, err)

	filepath, err := writeToml(bytes)
	assert.Nil(t, err)

	_, err = ReadConfigToml(filepath)
	assert.NotNil(t, err)
}

func TestMissingWebserverPort(t *testing.T) {
	missing_webserver_port := struct {
		DatabaseAddress  string `toml:"database_address"`
		DatabaseUser     string `toml:"database_user"`
		DatabasePassword string `toml:"database_password"`
		DatabaseName     string `toml:"database_name"`
		DatabasePort     string `toml:"database_port"`
	}{
		DatabaseAddress:  "test_database_address",
		DatabaseUser:     "test_database_user",
		DatabasePassword: "test_database_password",
		DatabaseName:     "test_database_name",
		DatabasePort:     "666",
	}

	bytes, err := toml.Marshal(missing_webserver_port)
	assert.Nil(t, err)

	filepath, err := writeToml(bytes)
	assert.Nil(t, err)

	_, err = ReadConfigToml(filepath)
	assert.NotNil(t, err)
}

func TestMissingDatabasePort(t *testing.T) {
	missing_database_address := struct {
		DatabaseAddress  string `toml:"database_address"`
		DatabaseUser     string `toml:"database_user"`
		DatabasePassword string `toml:"database_password"`
		DatabaseName     string `toml:"database_name"`
		WebserverPort    string `toml:"webserver_port"`
	}{
		DatabaseAddress:  "test_database_address",
		DatabaseUser:     "test_database_user",
		DatabasePassword: "test_database_password",
		DatabaseName:     "test_database_name",
		WebserverPort:    "99999",
	}

	bytes, err := toml.Marshal(missing_database_address)
	assert.Nil(t, err)

	filepath, err := writeToml(bytes)
	assert.Nil(t, err)

	_, err = ReadConfigToml(filepath)
	assert.NotNil(t, err)
}

func TestWrongDatabasePortValueType(t *testing.T) {
	missing_database_address := struct {
		DatabaseAddress  string `toml:"database_address"`
		DatabaseUser     string `toml:"database_user"`
		DatabasePassword string `toml:"database_password"`
		DatabaseName     string `toml:"database_name"`
		DatabasePort     string `toml:"database_port"`
		WebserverPort    int    `toml:"webserver_port"`
	}{
		DatabaseAddress:  "test_database_address",
		DatabaseUser:     "test_database_user",
		DatabasePassword: "test_database_password",
		DatabaseName:     "test_database_name",
		DatabasePort:     "String Should Not Work",
		WebserverPort:    99999,
	}

	bytes, err := toml.Marshal(missing_database_address)
	assert.Nil(t, err)

	filepath, err := writeToml(bytes)
	assert.Nil(t, err)

	_, err = ReadConfigToml(filepath)
	assert.NotNil(t, err)
}

func TestWrongwebserverPortValueType(t *testing.T) {
	missing_database_address := struct {
		DatabaseAddress  string `toml:"database_address"`
		DatabaseUser     string `toml:"database_user"`
		DatabasePassword string `toml:"database_password"`
		DatabaseName     string `toml:"database_name"`
		DatabasePort     int    `toml:"database_port"`
		WebserverPort    string `toml:"webserver_port"`
	}{
		DatabaseAddress:  "test_database_address",
		DatabaseUser:     "test_database_user",
		DatabasePassword: "test_database_password",
		DatabaseName:     "test_database_name",
		DatabasePort:     10,
		WebserverPort:    "99999",
	}

	bytes, err := toml.Marshal(missing_database_address)
	assert.Nil(t, err)

	filepath, err := writeToml(bytes)
	assert.Nil(t, err)

	_, err = ReadConfigToml(filepath)
	assert.NotNil(t, err)
}
