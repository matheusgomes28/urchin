package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/database"
	"github.com/rs/zerolog/log"
)

type AppSettings struct {
    Database_address string
    Database_port  int
	Database_user string
	Database_password string
}

func Run(app_settings AppSettings, database database.Database) error {
	r := gin.Default()
	r.MaxMultipartMemory = 1
	//r.LoadHTMLFiles("./templates/contact/contact-success.html", "./templates/contact/contact-failure.html")
	r.LoadHTMLGlob("templates/**/*")

	r.GET("/", makeHomeHandler(app_settings, database))

	// Contact form related endpoints
	r.GET("/contact", makeContactPageHandler(app_settings, database))
	r.POST("/contact-send", makeContactFormHandler())

	r.Static("/static", "./static")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	return nil
}

/// This function will act as the handler for
/// the home page
func makeHomeHandler(settings AppSettings, db database.Database) func(*gin.Context) {
	return func(c *gin.Context){
		posts, err := db.GetPosts()
		if err != nil {
			log.Error().Msgf("error loading posts: %v\n", err)
			return
		}

		c.HTML(http.StatusAccepted, "index", gin.H{"posts": posts})
	}
}
