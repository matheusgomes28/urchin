package index_test

import (
	"context"
	_ "database/sql"
	"embed"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
)

//go:generate ../../../migrations ./migrations

//go:embed migrations/*.sql
var embedMigrations embed.FS

var current_port int = 0

func runDatabaseServer(app_settings common.AppSettings) {
	pro := createTestDatabase(app_settings.DatabaseName)
	engine := sqle.NewDefault(pro)
	engine.Analyzer.Catalog.MySQLDb.AddRootAccount()

	session := memory.NewSession(sql.NewBaseSession(), pro)
	ctx := sql.NewContext(context.Background(), sql.WithSession(session))
	ctx.SetCurrentDatabase(app_settings.DatabaseName)

	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%d", app_settings.DatabaseAddress, app_settings.DatabasePort),
	}
	s, err := server.NewServer(config, engine, memory.NewSessionBuilder(pro), nil)
	if err != nil {
		panic(err)
	}
	if err = s.Start(); err != nil {
		panic(err)
	}
}

func createTestDatabase(name string) *memory.DbProvider {
	db := memory.NewDatabase(name)
	db.BaseDatabase.EnablePrimaryKeyIndexes()

	pro := memory.NewDBProvider(db)
	return pro
}

func waitForDb(app_settings common.AppSettings) (database.SqlDatabase, error) {

	for range 400 {
		database, err := database.MakeSqlConnection(
			app_settings.DatabaseUser,
			app_settings.DatabasePassword,
			app_settings.DatabaseAddress,
			app_settings.DatabasePort,
			app_settings.DatabaseName,
		)

		if err == nil {
			return database, nil
		}

		time.Sleep(25 * time.Millisecond)
	}

	return database.SqlDatabase{}, fmt.Errorf("database did not start")
}

func getAppSettings() common.AppSettings {
	app_settings := common.AppSettings{
		DatabaseAddress:  "localhost",
		DatabasePort:     3336 + current_port, // Initial port
		DatabaseUser:     "root",
		DatabasePassword: "",
		DatabaseName:     "urchin",
		WebserverPort:    8080,
	}
	current_port++

	return app_settings
}
