package admin_app

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"

	"github.com/fossoreslp/go-uuid-v4"
	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/rs/zerolog/log"
)

// TODO : need these endpoints
// r.GET("/images/:id", getImageHandler(&database))
// r.POST("/images", postImageHandler(&database))
// r.DELETE("/images", deleteImageHandler(&database))

func getImageHandler(database database.Database) func(*gin.Context) {
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

		c.JSON(http.StatusOK, image)
	}
}

func postImageHandler(app_settings common.AppSettings, database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10*1000000)
		form, err := c.MultipartForm()
		if err != nil {
			log.Error().Msgf("could not create multipart form: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("request type must be `multipart-form`", err))
			return
		}

		alt_text_array := form.Value["alt"]
		alt_text := "unknown"
		if len(alt_text_array) > 0 {
			alt_text = alt_text_array[0]
		}

		// Begin saving the file to the filesystem
		file_array := form.File["file"]
		if len(file_array) == 0 || file_array[0] == nil {
			log.Error().Msgf("could not get the file array: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no file provided for image upload",
			})
			return
		}

		file := file_array[0]
		allowed_types := []string{"image/jpeg", "image/png", "image/gif"}
		file_content_type := file.Header.Get("content-type")
		if !slices.Contains(allowed_types, file_content_type) {
			log.Error().Msgf("file type not supported")
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("file type not supported"))
			return
		}

		detected_content_type, err := getContentType(file)
		if err != nil || detected_content_type != file_content_type {
			log.Error().Msgf("the provided file does not match the provided content type")
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("provided file content is not allowed"))
			return
		}

		uuid, err := uuid.New()
		if err != nil {
			log.Error().Msgf("could not create the UUID: %v", err)
			c.JSON(http.StatusInternalServerError, common.ErrorRes("cannot create unique identifier", err))
			return
		}

		allowed_extensions := []string{"jpeg", "jpg", "png"}
		ext := filepath.Ext(file.Filename)
		// check ext is supported
		if ext == "" && slices.Contains(allowed_extensions, ext) {
			log.Error().Msgf("file extension is not supported %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("file extension is not supported", err))
			return
		}

		filename := fmt.Sprintf("%s%s", uuid.String(), ext)
		image_path := filepath.Join(app_settings.ImageDirectory, filename)
		err = c.SaveUploadedFile(file, image_path)
		if err != nil {
			log.Error().Msgf("could not save file: %v", err)
			c.JSON(http.StatusInternalServerError, common.ErrorRes("failed to upload image", err))
			return
		}
		// End saving to filesystem

		// Save metadata into the DB
		err = database.AddImage(uuid.String(), file.Filename, alt_text, ext)
		if err != nil {
			log.Error().Msgf("could not add image metadata to db: %v", err)
			os_err := os.Remove(image_path)
			if os_err != nil {
				log.Error().Msgf("could not remove image: %v", err)
				err = errors.Join(err, os_err)
			}
			c.JSON(http.StatusInternalServerError, common.ErrorRes("failed to save image", err))
			return
		}

		c.JSON(http.StatusOK, ImageIdResponse{
			Id: uuid.String(),
		})
	}
}

func deleteImageHandler(app_settings common.AppSettings, database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var delete_image_binding DeleteImageBinding
		err := c.ShouldBindUri(&delete_image_binding)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorRes("no id provided to delete image", err))
			return
		}

		image, err := database.GetImage(delete_image_binding.Id)
		if err != nil {
			log.Error().Msgf("failed to delete image: %v", err)
			c.JSON(http.StatusNotFound, common.ErrorRes("could not delete image", err))
			return
		}

		filename := fmt.Sprintf("%s%s", image.Uuid, image.Ext)
		image_path := filepath.Join(app_settings.ImageDirectory, filename)
		err = os.Remove(image_path)
		if err != nil {
			log.Warn().Msgf("could not delete stored image file: %v", err)
			// No return because we have to remove the database entry nonetheless.
		}

		err = database.DeleteImage(delete_image_binding.Id)
		if err != nil {
			log.Error().Msgf("failed to delete image with id %s: %v", delete_image_binding.Id, err)
			c.JSON(http.StatusNotFound, common.ErrorRes("could not delete image", err))
			return
		}

		c.JSON(http.StatusOK, ImageIdResponse{
			delete_image_binding.Id,
		})
	}
}

func getContentType(file_header *multipart.FileHeader) (string, error) {
	// Check if the content matches the provided type.
	image_file, err := file_header.Open()
	if err != nil {
		log.Error().Msgf("could not open file for check.")
		return "", err
	}

	// According to the documentation only the first `512` bytes are required for verifying the content type
	tmp_buffer := make([]byte, 512)
	_, read_err := image_file.Read(tmp_buffer)
	if read_err != nil {
		log.Error().Msgf("could not read into temp buffer")
		return "", read_err
	}
	return getContentTypeFromData(tmp_buffer), nil
}

func getContentTypeFromData(data []byte) string {
	return http.DetectContentType(data[:512])
}
