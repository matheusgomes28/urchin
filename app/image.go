package app

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"slices"
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

	// Get all the files inside the image directory
	files, err := os.ReadDir(app_settings.ImageDirectory)
	if err != nil {
		log.Error().Msgf("could not read files in image directory: %v", err)
		return []byte{}, err
	}

	// Filter all the non-images out of the list
	valid_images := make([]common.Image, 0)
	valid_extensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	for n, file := range files {

		// TODO : This is surely not the best way
		//        to implement pagination in for loops
		if n >= limit {
			break
		}

		if n < offset {
			continue
		}

		filename := file.Name()
		ext := path.Ext(file.Name())
		if slices.Contains(valid_extensions, ext) {

			image := common.Image{
				Uuid:    filename[:len(filename)-len(ext)],
				Name:    filename,
				AltText: "undefined", // TODO : perhaps remove this
				Ext:     ext,
			}
			valid_images = append(valid_images, image)
		}
	}

	index_view := views.MakeImagesPage(valid_images, app_settings.AppNavbar.Links)
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

	// if not cached, create the cache
	filename := get_image_binding.Filename
	ext := path.Ext(get_image_binding.Filename)
	name := filename[:len(filename)-len(ext)]

	image := common.Image{
		Uuid:    name,
		Name:    filename,
		AltText: "undefined",
		Ext:     ext,
	}
	index_view := views.MakeImagePage(image, app_settings.AppNavbar.Links)
	html_buffer := bytes.NewBuffer(nil)

	err := index_view.Render(c, html_buffer)
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

		// TODO : How do we get the filename?
		image_path := filepath.Join(app_settings.ImageDirectory, get_image_binding.Filename)
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
