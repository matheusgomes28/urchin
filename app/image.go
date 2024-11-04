package app

import (
	"bytes"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/matheusgomes28/urchin/views"
	"github.com/rs/zerolog/log"
)

// Since there are no builtin sets in go, we are using a map to improve the performance when checking for valid extensions
// by creating a map with the valid extensions as keys and using an existence check.
var valid_extensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

func imagesHandler(c *gin.Context, app_settings common.AppSettings, database database.Database) ([]byte, error) {
	// TODO: Implement rendering.
	pageNum := 1 // Default to page 0
	if pageNumQuery := c.Param("num"); pageNumQuery != "" {
		num, err := strconv.Atoi(pageNumQuery)
		if err == nil && num > 0 {
			pageNum = num
		} else {
			log.Error().Msgf("Invalid page number: %s", pageNumQuery)
		}
	}

	limit := 10 // or whatever limit you want
	offset := (pageNum - 1) * limit

	// Get all the files inside the image directory
	files, err := os.ReadDir(app_settings.ImageDirectory)
	if err != nil {
		log.Error().Msgf("could not read files in image directory: %v", err)
		return []byte{}, err
	}

	// Filter all the non-images out of the list
	valid_images := make([]common.Image, 0)
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
		// Checking for the existence of a value in a map takes O(1) and therefore it's faster than
		// iterating over a string slice
		_, ok := valid_extensions[ext]
		if !ok {
			continue
		}

		image := common.Image{
			Uuid: filename[:len(filename)-len(ext)],
			Name: filename,
			Ext:  ext,
		}
		valid_images = append(valid_images, image)
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
		Uuid: name,
		Name: filename,
		Ext:  ext,
	}

	return renderHtml(c, views.MakeImagePage(image, app_settings.AppNavbar.Links))
}
