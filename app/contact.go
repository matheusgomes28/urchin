package app

import (
	"bytes"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/matheusgomes28/urchin/views"
	"github.com/rs/zerolog/log"
)

func makeContactFormHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		if err := c.Request.ParseForm(); err != nil {
			log.Error().Msgf("could not parse form %v", err)
			if err = render(c, http.StatusOK, views.MakeContactFailure("unknown", err.Error())); err != nil {
				log.Error().Msgf("could not render %v", err)
			}
			return
		}

		email := c.Request.FormValue("email")
		name := c.Request.FormValue("name")
		message := c.Request.FormValue("message")

		// Parse email
		_, err := mail.ParseAddress(email)
		if err != nil {
			log.Error().Msgf("could not parse email: %v", err)
			if err = render(c, http.StatusOK, views.MakeContactFailure(email, err.Error())); err != nil {
				log.Error().Msgf("could not render: %v", err)
			}
			return
		}

		// Make sure name and message is reasonable
		if len(name) > 200 {
			if err = render(c, http.StatusOK, views.MakeContactFailure(email, "name too long (200 chars max)")); err != nil {
				log.Error().Msgf("could not render: %v", err)
			}
			return
		}

		if len(message) > 10000 {
			if err = render(c, http.StatusOK, views.MakeContactFailure(email, "message too long (1000 chars max)")); err != nil {
				log.Error().Msgf("could not render: %v", err)
			}
			return
		}

		if err = render(c, http.StatusOK, views.MakeContactSuccess(email, name)); err != nil {
			log.Error().Msgf("could not render: %v", err)
		}
	}
}

// TODO : This is a duplicate of the index handler... abstract
func contactHandler(c *gin.Context, app_settings common.AppSettings, db *database.Database) ([]byte, error) {
	index_view := views.MakeContactPage()
	html_buffer := bytes.NewBuffer(nil)
	if err := index_view.Render(c, html_buffer); err != nil {
		log.Error().Msgf("could not render: %v", err)
	}

	return html_buffer.Bytes(), nil
}
