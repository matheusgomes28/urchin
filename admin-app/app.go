package admin_app

import (
	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
)

type PostBinding struct {
	Id string `uri:"id" binding:"required"`
}

type AddPostRequest struct {
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Content string `json:"content"`
}

type ChangePostRequest struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Content string `json:"content"`
}

type DeletePostRequest struct {
	Id int `json:"id"`
}

func SetupRoutes(app_settings common.AppSettings, database database.Database) *gin.Engine {

	r := gin.Default()
	r.MaxMultipartMemory = 1

	r.GET("/posts/:id", getPostHandler(database))
	r.POST("/posts", postPostHandler(database))
	r.PUT("/posts", putPostHandler(database))
	r.DELETE("/posts", deletePostHandler(database))

	// CRUD images
	// r.GET("/images/:id", getImageHandler(&database))
	r.POST("/images", postImageHandler(app_settings, database))
	// r.DELETE("/images", deleteImageHandler(&database))

	return r
}
