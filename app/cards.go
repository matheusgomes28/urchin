package app

import (
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/matheusgomes28/urchin/views"
	"github.com/rs/zerolog/log"
)

type CardIdBinding struct {
	Id string `uri:"id" binding:"required"`
}

func cardHandler(c *gin.Context, app_settings common.AppSettings, database database.Database) ([]byte, error) {
	var card_binding CardIdBinding
	if err := c.ShouldBindUri(&card_binding); err != nil {
		return nil, err
	}

	// Get the post with the ID
	card, err := database.GetCard(card_binding.Id)
	if err != nil {
		return nil, err
	}

	// Generate HTML page
	post_view := views.MakeCardPage(card.JsonData)
	html_buffer := bytes.NewBuffer(nil)
	err = post_view.Render(c, html_buffer)
	if err != nil {
		log.Error().Msgf("%s", err)
	}

	return html_buffer.Bytes(), nil
}
