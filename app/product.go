package app

import (
	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/matheusgomes28/urchin/views"
)

func productHandler(c *gin.Context, app_settings common.AppSettings, db database.Database) ([]byte, error) {
	return renderHtml(c, views.MakeProductPage(app_settings.AppNavbar.Links))
}