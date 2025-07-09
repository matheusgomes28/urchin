package admin_app

import (
	"net/http"
	"strings"

	_ "github.com/matheusgomes28/urchin/docs/admin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
)

func SetupRoutes(app_settings common.AppSettings, database database.Database) *gin.Engine {

	r := gin.Default()
	r.MaxMultipartMemory = 1

	r.GET("/posts/:id", getPostHandler(database))
	r.POST("/posts", postPostHandler(database))
	r.PUT("/posts", putPostHandler(database))
	r.DELETE("/posts/:id", deletePostHandler(database))

	r.POST("/images", postImageHandler(app_settings))
	r.DELETE("/images/:name", deleteImageHandler(app_settings))

	r.POST("/pages", postPageHandler(database))

	// For container health purposes
	r.Any("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, PostIdResponse{Id: 0})
	})

	// Swagger routes
	r.GET("/swagger/*any", func(c *gin.Context) {
		// Don't serve index.html again if it's already handled above
		if strings.HasSuffix(c.Request.URL.Path, "/swagger/") {
			c.Request.RequestURI = "/swagger/index.html"
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})

	return r
}
