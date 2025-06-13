package app

import (
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
)

type GalleryBinding struct {
	Name string `uri:"name" binding:"required"`
}

func galleryHandler(c *gin.Context, app_settings common.AppSettings, database database.Database) ([]byte, error) {
	var get_gallery_binding struct {
		Name string `uri:"name" binding:"required"`
	}

	if err := c.ShouldBindUri(&get_gallery_binding); err != nil {
		return nil, err
	}

	html_buffer := bytes.NewBuffer(nil)
	return html_buffer.Bytes(), nil
}
