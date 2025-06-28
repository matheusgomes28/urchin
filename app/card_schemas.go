package app

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/matheusgomes28/urchin/views"
	"github.com/rs/zerolog/log"
)

func getSchemasHandler(c *gin.Context, app_settings common.AppSettings, database database.Database) ([]byte, error) {
	// Lee offset y limit de la query (?offset=0&limit=10)
	// Un valor de 0 para el límite significa "sin límite".
	offsetStr := c.Param("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}
	limitStr := c.Param("limit")
	if limitStr == "" {
		limitStr = "0"
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
		return []byte{}, err
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return []byte{}, err
	}

	// Si no se especifica un límite, se obtienen todos los posts.
	// Si se especifica, se usa para la paginación.
	schemas, err := database.GetCardSchemas(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return []byte{}, err
	}
	schemas_view := views.MakeAllSchemas(schemas, app_settings.AppNavbar.Links)
	html_buffer := bytes.NewBuffer(nil)

	err = schemas_view.Render(c, html_buffer)
	if err != nil {
		log.Error().Msgf("Could not render index: %v", err)
		return []byte{}, err
	}

	return html_buffer.Bytes(), nil
}
