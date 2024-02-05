package admin_app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fossoreslp/go-uuid-v4"
	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/rs/zerolog/log"
)

type AddImageRequest struct {
	Alt string `json:"alt"`
}

// r.GET("/images/:id", getImageHandler(&database))
// r.POST("/images", postImageHandler(&database))
// r.DELETE("/images", deleteImageHandler(&database))

func getImageHandler(database *database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		// Get the image from database
	}
}

func postImageHandler(app_settings common.AppSettings, database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			log.Error().Msgf("could nto create multipart form: %v", err)
			return
		}

		alt_text_array := form.Value["alt"]
		alt_text := "unknown"
		if len(alt_text_array) > 0 {
			alt_text = alt_text_array[0]
		}

		// Begin saving the file to the filesystem
		file_array := form.File["file"]
		if len(file_array) == 0 {
			log.Error().Msgf("could not get the file array: %v", err)
			return
		}
		file := file_array[0]
		if file == nil {
			log.Error().Msgf("could not upload file: %v", err)
			return
		}

		uuid, err := uuid.New()
		if err != nil {
			log.Error().Msgf("could not create the UUID: %v", err)
			return
		}

		ext := filepath.Ext(file.Filename)
		// check ext is supported
		if ext == "" {
			log.Error().Msgf("could not get file extension %v", err)
			return
		}

		filename := fmt.Sprintf("%s.%s", uuid.String(), ext)
		image_path := filepath.Join(app_settings.ImageDirectory, filename)
		err = c.SaveUploadedFile(file, image_path)
		if err != nil {
			log.Error().Msgf("could not save file: %v", err)
			return
		}
		// End saving to filesystem

		// Save metadata into the DB
		err = database.AddImage(uuid.String(), file.Filename, alt_text)
		if err != nil {
			log.Error().Msgf("could not add image metadata to db: %v", err)
			err := os.Remove(image_path)
			if err != nil {
				log.Error().Msgf("could not remove image: %v", err)
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": uuid.String(),
		})
	}
}
