package admin_app

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/rs/zerolog/log"
)

func postCardHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var add_card_request AddCardRequest
		if c.Request.Body == nil {
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("no request body provided"))
			return
		}
		decoder := json.NewDecoder(c.Request.Body)
		err := decoder.Decode(&add_card_request)

		if err != nil {
			log.Warn().Msgf("invalid post request: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("invalid request body", err))
			return
		}

		// TODO : Sanity checks that everything inside the
		// TODO : request makes sense. I.e. content is json,
		// TODO : i.e json content matches the schema, etc.
		// err = checkRequiredData(add_card_request)
		// if err != nil {
		// 	log.Error().Msgf("failed to add post required data is missing: %v", err)
		// 	c.JSON(http.StatusBadRequest, common.ErrorRes("missing required data", err))
		// 	return
		// }

		id, err := database.AddCard(
			add_card_request.Title,
			add_card_request.Image,
			add_card_request.Schema,
			add_card_request.Content,
		)
		if err != nil {
			log.Error().Msgf("failed to add card: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not add post", err))
			return
		}

		c.JSON(http.StatusOK, CardIdResponse{
			id,
		})
	}
}
