package app

import (
	"bytes"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/common"
	"github.com/matheusgomes28/database"
	"github.com/matheusgomes28/views"
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
			render(c, http.StatusOK, views.MakeContactFailure(email, err.Error()))
			return
		}

		// Make sure name and message is reasonable
		if len(name) > 200 {
			render(c, http.StatusOK, views.MakeContactFailure(email, "name too long (200 chars max)"))
			return
		} 

		if len(message) > 10000 {
			render(c, http.StatusOK, views.MakeContactFailure(email, "message too long (1000 chars max)"))
			return
		}

		render(c, http.StatusOK, views.MakeContactSuccess(email, name))
	}
}

// TODO : This is a duplicate of the index handler... abstract
func contactHandler(c *gin.Context, app_settings common.AppSettings, db *database.Database) ([]byte, error) {
	index_view := views.MakeContactPage()
	html_buffer := bytes.NewBuffer(nil)
	index_view.Render(c, html_buffer)

	return html_buffer.Bytes(), nil
}
