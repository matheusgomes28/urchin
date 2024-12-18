package admin_app

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/rs/zerolog/log"
)

func postSchemaHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var add_schema_request AddCardSchemaRequest
		if c.Request.Body == nil {
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("no request body provided"))
			return
		}

		// Validate the content of the schema
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("could not read request body"))
			return
		}

		if !json.Valid(body) {
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("invalid json given in request body"))
			return
		}

		err = json.Unmarshal(body, &add_schema_request)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("could not unmarshal json into request"))
			return
		}

		if !checkSchemaValues(add_schema_request) {
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("invalid schema inputs"))
			return
		}

		id, err := database.AddCardSchema(
			add_schema_request.JsonId,
			add_schema_request.JsonSchema,
			add_schema_request.JsonTitle,
			add_schema_request.Schema,
		)
		if err != nil {
			log.Error().Msgf("failed to add card schema: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not add card schema", err))
			return
		}

		c.JSON(http.StatusOK, CardIdResponse{
			id,
		})
	}
}

func checkSchemaValues(add_schema_request AddCardSchemaRequest) bool {

	if add_schema_request.JsonId == "" {
		return false
	}
	if add_schema_request.JsonSchema == "" {
		return false
	}
	if add_schema_request.JsonTitle == "" {
		return false
	}
	if add_schema_request.Schema == "" {
		return false
	}

	return true
}
