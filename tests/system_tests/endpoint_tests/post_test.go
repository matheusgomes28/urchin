package index_test

import (
	_ "database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	"github.com/matheusgomes28/urchin/app"
	"github.com/pressly/goose/v3"
)

func TestPostExists(t *testing.T) {

	// This is gonna be the in-memory mysql
	app_settings := getAppSettings()
	go runDatabaseServer(app_settings)
	database, err := waitForDb(app_settings)
	require.Nil(t, err)
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("mysql"); err != nil {
		require.Nil(t, err)
	}

	if err := goose.Up(database.Connection, "migrations"); err != nil {
		require.Nil(t, err)
	}

	r := app.SetupRoutes(app_settings, database)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/post/1", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}
