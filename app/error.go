package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/views"
)

func NotFoundHandler(app_settings common.AppSettings) func(*gin.Context) {
	handler := func(c *gin.Context) {
		buffer, err := renderHtml(c, views.MakeNotFoundPage(app_settings.AppNavbar.Links, app_settings.AppNavbar.Dropdowns))
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrorRes("could not render HTML", err))
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", buffer)
	}

	return handler
}

func ErrorHandler(error_msg string, app_settings common.AppSettings) func(*gin.Context) {
	handler := func(c *gin.Context) {
		buffer, err := renderHtml(c, views.MakeNotFoundPage(app_settings.AppNavbar.Links, app_settings.AppNavbar.Dropdowns))
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrorRes("could not render HTML", err))
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", buffer)
	}

	return handler
}
