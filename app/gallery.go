package app

import (
	"bytes"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/matheusgomes28/urchin/views"
)

// This gets the images from all the declared
// manifest files in the gallery settings.
//
// The manifest files in "gallery.links" should
// be relative to the image directory
func getGalleryImages(gallery common.Gallery, app_settings common.AppSettings) ([]common.Image, error) {
	images, err := common.GetImages(gallery.Images, len(gallery.Images), 1, app_settings)
	if err != nil {
		return []common.Image{}, err
	}

	return images, nil
}

func galleryHandler(c *gin.Context, app_settings common.AppSettings, database database.Database) ([]byte, error) {
	var get_gallery_binding struct {
		Name string `uri:"name" binding:"required"`
	}

	if err := c.ShouldBindUri(&get_gallery_binding); err != nil {
		return []byte{}, err
	}

	gallery, exists := app_settings.Galleries[get_gallery_binding.Name]
	if !exists {
		return []byte{}, fmt.Errorf("requested gallery `%s` does not exist", gallery.Name)
	}

	// TODO : Get valid images for a gallery
	images, err := getGalleryImages(gallery, app_settings)
	if err != nil {
		return []byte{}, fmt.Errorf("could not get gallery: %v", err)
	}

	gallery_view := views.MakeImagesPage(images, app_settings.AppNavbar.Links, app_settings.AppNavbar.Dropdowns)
	html_buffer := bytes.NewBuffer(nil)
	err = gallery_view.Render(c, html_buffer)
	if err != nil {
		return []byte{}, err
	}

	return html_buffer.Bytes(), nil
}
