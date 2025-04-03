package endpoint_tests

import (
	_ "database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/matheusgomes28/urchin/app"
	"github.com/matheusgomes28/urchin/tests/helpers"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
)

func TestPostExists(t *testing.T) {

	// This is gonna be the in-memory mysql
	app_settings := helpers.GetAppSettings()
	database, err := helpers.WaitForDb(app_settings)
	require.Nil(t, err)
	goose.SetBaseFS(helpers.EmbedMigrations)

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
