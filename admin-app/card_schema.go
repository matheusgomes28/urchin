package admin_app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaptinlin/jsonschema"
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
			error_msg := fmt.Errorf("could not unmarshall json request: %v", err)
			c.JSON(http.StatusBadRequest, common.MsgErrorRes(error_msg.Error()))
			return
		}

		if err = checkSchemaValues(add_schema_request); err != nil {
			c.JSON(http.StatusBadRequest, common.MsgErrorRes(err.Error()))
			return
		}

		id, err := database.AddCardSchema(
			add_schema_request.JsonSchema,
			add_schema_request.JsonTitle,
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

func checkSchemaValues(add_schema_request AddCardSchemaRequest) error {

	if add_schema_request.JsonSchema == "" {
		return fmt.Errorf("`schema` cannot be empty")
	}
	if add_schema_request.JsonTitle == "" {
		return fmt.Errorf("`title` cannot be empty")
	}

	schema_compiler := jsonschema.NewCompiler()
	_, err := schema_compiler.Compile([]byte(add_schema_request.JsonSchema))

	if err != nil {
		return fmt.Errorf("`schema` is invalid: %v", err)
	}

	return nil
}
