package admin_app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaptinlin/jsonschema"
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

		// Check that the schema exists
		schema, err := database.GetCardSchema(add_card_request.Schema)
		if err != nil {
			log.Error().Msgf("card schema does not exist: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("card schema does not exist", err))
			return
		}

		err = validateCardAgainstSchema(add_card_request.Content, schema.Schema)
		if err != nil {
			log.Error().Msgf(err.Error())
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not add card", err))
			return
		}

		id, err := database.AddCard(
			add_card_request.Image,
			add_card_request.Schema,
			add_card_request.Content,
		)
		if err != nil {
			log.Error().Msgf("failed to add card: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not add card", err))
			return
		}

		c.JSON(http.StatusOK, CardIdResponse{
			id,
		})
	}
}

func validateCardAgainstSchema(card_data string, json_schema string) error {

	// Parse the schema here
	schema_compiler := jsonschema.NewCompiler()
	schema, err := schema_compiler.Compile([]byte(json_schema))

	if err != nil {
		return fmt.Errorf("failed to compile the json_schema from db: %v", err)
	}

	json_map := make(map[string]interface{})
	err = json.Unmarshal([]byte(card_data), &json_map)
	if err != nil {
		return fmt.Errorf("failed to parse card json : %v", err)
	}

	result := schema.Validate(json_map)
	if !result.IsValid() {
		details, _ := json.MarshalIndent(result.ToList(), "", "  ")
		return fmt.Errorf("failed to check vard data against schema: %v", string(details))
	}

	return nil
}
