package endpoint_tests

import (
	_ "database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	"github.com/matheusgomes28/urchin/app"
	"github.com/matheusgomes28/urchin/tests/helpers"
)

func TestIndexPagePing(t *testing.T) {
	// Usual database setup
	app_settings := helpers.GetAppSettings()
	cleanup, db, err := helpers.SetupDb(app_settings)
	require.Nil(t, err)
	defer func() { require.Nil(t, cleanup()) }()

	r := app.SetupRoutes(app_settings, db)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

// TODO : Uncomment this
func TestIndexPagePostExists(t *testing.T) {

	// Usual db setup
	app_settings := helpers.GetAppSettings()
	cleanup, db, err := helpers.SetupDb(app_settings)
	require.Nil(t, err)
	defer func() { require.Nil(t, cleanup()) }()

	r := app.SetupRoutes(app_settings, db)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	require.Contains(t, w.Body.String(), "My Very First Post")
}
