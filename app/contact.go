package app

import (
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/database"
	"github.com/rs/zerolog/log"
)

func makeContactFormHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		c.Request.ParseForm()
		email := c.Request.FormValue("email")
		name := c.Request.FormValue("name")
		message := c.Request.FormValue("message")

		// Parse email
		_, err := mail.ParseAddress(email)
		if err != nil { 	
			c.HTML(http.StatusOK, "contact-failure.html", gin.H{
				"email": email,
				"error": "invalid email",
			})
			return
		}

		// Make sure name and message is reasonable
		if len(name) > 200 {
			c.HTML(http.StatusOK, "contact-failure.html", gin.H{
				"email": email,
				"error": "invalid name",
			})
			return
		} 

		if len(message) > 10000 {
			c.HTML(http.StatusOK, "contact-failure.html", gin.H{
				"email": email,
				"error": "message too big",
			})
			return
		}

		c.HTML(http.StatusOK, "contact-success.html", gin.H{
			"name": name,
			"email": email,
		})
	}
}

// TODO : This is a duplicate of the index handler... abstract
func makeContactPageHandler(settings AppSettings, db database.Database) func(*gin.Context) {
	return func(c *gin.Context){
		posts, err := db.GetPosts()
		if err != nil {
			log.Error().Msgf("error loading posts: %v\n", err)
			return
		}

		c.HTML(http.StatusAccepted, "contact", gin.H{"posts": posts})
	}
}
