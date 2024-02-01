package app

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/common"
	"github.com/matheusgomes28/database"
	"github.com/matheusgomes28/views"
	"github.com/matheusgomes28/views/tailwind"
	"github.com/rs/zerolog/log"
)

func Run(app_settings common.AppSettings, database database.Database) error {
	r := gin.Default()
	r.MaxMultipartMemory = 1

	r.GET("/", makeHomeHandler(app_settings, database, views.MakeIndex))
	r.GET("/tailwind", makeHomeHandler(app_settings, database, tailwind.MakeIndex))

	// Contact form related endpoints
	r.GET("/contact", makeContactPageHandler(app_settings, database))
	r.POST("/contact-send", makeContactFormHandler())


	// Post related endpoints
	r.GET("/post/:id", makePostHandler(database))

	r.Static("/static", "./static")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	return nil
}

/// This function will act as the handler for
/// the home page
func makeHomeHandler(settings common.AppSettings, db database.Database, factory func([]common.Post) templ.Component) func(*gin.Context) {
	return func(c *gin.Context){
		posts, err := db.GetPosts()
		if err != nil {
			log.Error().Msgf("error loading posts: %v\n", err)
			return
		}

		render(c, http.StatusOK, factory(posts))
	}
}
