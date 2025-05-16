package endpoint_tests

import (
	_ "database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/matheusgomes28/urchin/app"
	"github.com/matheusgomes28/urchin/tests/helpers"
	"github.com/stretchr/testify/require"
)

func TestPostExists(t *testing.T) {
	// Usual db setup
	app_settings := helpers.GetAppSettings()
	cleanup, db, err := helpers.SetupDb(app_settings)
	require.Nil(t, err)
	defer func() { require.Nil(t, cleanup()) }()

	r := app.SetupRoutes(app_settings, db)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/post/1", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}
