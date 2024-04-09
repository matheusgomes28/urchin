package app

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/matheusgomes28/urchin/views"
	"github.com/rs/zerolog/log"
)

func imagesHandler(c *gin.Context, app_settings common.AppSettings, database database.Database) ([]byte, error) {
	// TODO: Implement rendering.
	pageNum := 0 // Default to page 0
	if pageNumQuery := c.Param("num"); pageNumQuery != "" {
		num, err := strconv.Atoi(pageNumQuery)
		if err == nil && num > 0 {
			pageNum = num
		} else {
			log.Error().Msgf("Invalid page number: %s", pageNumQuery)
		}
	}

	limit := 10 // or whatever limit you want
	offset := max((pageNum-1)*limit, 0)
	images, err := database.GetImages(offset, limit)
	if err != nil {
		return nil, err
	}

	// if not cached, create the cache
	index_view := views.MakeImagesPage(images, app_settings.AppNavbar.Links)
	html_buffer := bytes.NewBuffer(nil)

	err = index_view.Render(c, html_buffer)
	if err != nil {
		log.Error().Msgf("Could not render index: %v", err)
		return []byte{}, err
	}

	return html_buffer.Bytes(), nil
}

func imageHandler(c *gin.Context, app_settings common.AppSettings, database database.Database) ([]byte, error) {
	var get_image_binding common.ImageIdBinding
	if err := c.ShouldBindUri(&get_image_binding); err != nil {
		return nil, err
	}

	image, err := database.GetImage(get_image_binding.Id)
	if err != nil {
		return nil, err
	}

	// if not cached, create the cache
	index_view := views.MakeImagePage(image, app_settings.AppNavbar.Links)
	html_buffer := bytes.NewBuffer(nil)

	err = index_view.Render(c, html_buffer)
	if err != nil {
		log.Error().Msgf("Could not render index: %v", err)
		return []byte{}, err
	}

	return html_buffer.Bytes(), nil
}

func getImageHandler(app_settings common.AppSettings, database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var get_image_binding common.ImageIdBinding
		if err := c.ShouldBindUri(&get_image_binding); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not get image id", err))
			return
		}

		image, err := database.GetImage(get_image_binding.Id)
		if err != nil {
			log.Error().Msgf("failed to get image: %v", err)
			c.JSON(http.StatusNotFound, common.ErrorRes("could not get image", err))
			return
		}

		filename := fmt.Sprintf("%s%s", image.Uuid, image.Ext)
		image_path := filepath.Join(app_settings.ImageDirectory, filename)
		file, err := os.Open(image_path)
		if err != nil {
			log.Error().Msgf("failed to load saved image: %v", err)
			c.JSON(http.StatusNotFound, common.ErrorRes("image not found", err))
			return
		}

		data, err := io.ReadAll(file)
		if err != nil {
			log.Error().Msgf("failed to load saved image: %v", err)
			c.JSON(http.StatusNotFound, common.ErrorRes("image not found", err))
			return
		}

		c.Data(http.StatusOK, getContentTypeFromData(data), data)
	}
}

func getContentTypeFromData(data []byte) string {
	return http.DetectContentType(data[:512])
}
