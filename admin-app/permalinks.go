package admin_app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/rs/zerolog/log"
)

// / postPageHandler is the function handling the endpoint for adding new pages
func postPermalinkHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		add_permalink_request := struct {
			Permalink string `uri:"permalink" binding:"required"`
			PostId    int    `uri:"post_id" binding:"required"`
		}{}

		err := c.ShouldBindUri(&add_permalink_request)
		if err != nil {
			log.Error().Msgf("invalid request for adding permalink: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not add permalink", err))
			return
		}

		permalinkDb := common.Permalink{
			Path:   add_permalink_request.Permalink,
			PostId: add_permalink_request.PostId,
		}

		id, err := database.AddPermalink(permalinkDb)

		if err != nil {
			log.Error().Msgf("failed to add post: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not add post", err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	}
}
