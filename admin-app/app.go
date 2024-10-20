package admin_app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/matheusgomes28/urchin/plugins"
	lua "github.com/yuin/gopher-lua"
)

func SetupRoutes(app_settings common.AppSettings, database database.Database, shortcode_handlers map[string]*lua.LState, hooks map[string]plugins.Hook) *gin.Engine {

	r := gin.Default()
	r.MaxMultipartMemory = 1

	post_hook, ok := hooks["add_post"]
	if !ok {
		log.Fatalf("could not find add_post hook")
	}

	r.GET("/posts/:id", getPostHandler(database))
	r.POST("/posts", postPostHandler(database, shortcode_handlers, post_hook.(plugins.PostHook)))
	r.PUT("/posts", putPostHandler(database))
	r.DELETE("/posts/:id", deletePostHandler(database))

	r.POST("/images", postImageHandler(app_settings))
	r.DELETE("/images/:name", deleteImageHandler(app_settings))

	r.POST("/pages", postPageHandler(database))

	return r
}
