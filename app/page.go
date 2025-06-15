package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/matheusgomes28/urchin/views"
)

func pageHandler(c *gin.Context, app_settings common.AppSettings, database database.Database) ([]byte, error) {
	var page_binding common.PageLinkBinding
	err := c.ShouldBindUri(&page_binding)

	if err != nil || len(page_binding.Link) == 0 {
		// TODO : we should be serving an error page
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page uri"})

		return nil, err
	}

	// Get the page with the ID
	page, err := database.GetPage(page_binding.Link)

	if err != nil || page.Content == "" {
		// TODO : serve the error page instead
		c.JSON(http.StatusNotFound, gin.H{"error": "Page Not Found"})
		return nil, err
	}

	// Generate HTML page
	page.Content = string(mdToHTML([]byte(page.Content)))

	return renderHtml(c, views.MakePage(page.Title, page.Content, app_settings.AppNavbar.Links, app_settings.AppNavbar.Dropdowns))
}
