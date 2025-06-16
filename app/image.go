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

func imagesHandler(c *gin.Context, app_settings common.AppSettings, database database.Database) ([]byte, error) {
	pageNum := 1 // Default to page 0
	if pageNumQuery := c.Param("num"); pageNumQuery != "" {
		num, err := strconv.Atoi(pageNumQuery)
		if err == nil && num > 0 {
			pageNum = num
		} else {
			log.Error().Msgf("Invalid page number: %s", pageNumQuery)
		}
	}

	// Get all the files inside the image directory
	files, err := os.ReadDir(app_settings.ImageDirectory)
	if err != nil {
		log.Error().Msgf("could not read files in image directory: %v", err)
		return []byte{}, err
	}

	filepaths := common.Map(files, func(file os.DirEntry) string {
		return file.Name()
	})
	filepaths = common.Filter(filepaths, func(filepath string) bool {
		ext := path.Ext(filepath)
		return ext == ".json"
	})

	valid_images, err := common.GetImages(filepaths, 10, pageNum, app_settings)
	if err != nil {
		return []byte{}, err
	}

	index_view := views.MakeImagesPage(valid_images, app_settings.AppNavbar.Links, app_settings.AppNavbar.Dropdowns)
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

	return renderHtml(c, views.MakeImagePage(image, app_settings.AppNavbar.Links, app_settings.AppNavbar.Dropdowns))
}
