package admin_app

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/fossoreslp/go-uuid-v4"
	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/draw"
)

var allowed_extensions = map[string]bool{
	".jpeg": true, ".jpg": true, ".png": true,
}

var allowed_content_types = map[string]bool{
	"image/jpeg": true, "image/png": true, "image/gif": true,
}

// Calculates the best suitable resize box for an image with the given
// `original_width`, `original_height`, such that the returned (width, height)
// dimensions will always be at most `max_width`, `max_height`
func calculateResizeBox(original_width, original_height, max_height, desired_height int) (int, int) {
	width_ratio := float64(max_height) / float64(original_width)
	new_height := int(float64(original_height) * width_ratio)
	if new_height <= desired_height {
		return max_height, new_height
	}

	height_ratio := float64(desired_height) / float64(original_height)
	new_width := int(float64(original_width) * height_ratio)
	return new_width, int(height_ratio)
}

func resizeImage(src_path, dst_path string, desired_width, desired_height int) error {
	// Open the source file
	file, err := os.Open(src_path)
	if err != nil {
		return fmt.Errorf("could not open source image: %v", err)
	}
	defer file.Close()

	// Decode image
	img, format, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("could not decode image: %v", err)
	}

	// Calculate new size
	bounds := img.Bounds()
	new_width, new_height := calculateResizeBox(bounds.Dx(), bounds.Dy(), desired_width, desired_height)
	dst := image.NewRGBA(image.Rect(0, 0, new_width, new_height))

	// Resize using golang.org/x/image/draw
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	// Create new file
	out, err := os.Create(dst_path)
	if err != nil {
		return fmt.Errorf("could not create output file: %v", err)
	}
	defer out.Close()

	// Save based on format
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(out, dst, &jpeg.Options{Quality: 85})
	case "png":
		err = png.Encode(out, dst)
	case "gif":
		err = png.Encode(out, dst)
		log.Warn().Msg("GIF detected: resized image saved as PNG format with .gif extension")
	default:
		return fmt.Errorf("unsupported image format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("could not encode resized image: %v", err)
	}

	return nil
}

func createMinifiedImages(image_path string) error {

	image_ext := path.Ext(image_path)
	if len(image_ext) == 0 {
		return fmt.Errorf("invalid image path: %s", image_path)
	}
	image_name := image_path[0 : len(image_path)-len(image_ext)]

	image_types := []struct {
		FileName  string
		MaxWidth  int
		MaxHeight int
	}{
		{FileName: fmt.Sprintf("%s_small%s", image_name, image_ext), MaxWidth: 200, MaxHeight: 200},
		{FileName: fmt.Sprintf("%s_medium%s", image_name, image_ext), MaxWidth: 400, MaxHeight: 400},
		{FileName: fmt.Sprintf("%s_large%s", image_name, image_ext), MaxWidth: 600, MaxHeight: 600},
	}

	for _, img_type := range image_types {
		err := resizeImage(image_path, img_type.FileName, img_type.MaxWidth, img_type.MaxHeight)
		if err != nil {
			return fmt.Errorf("could not resize image: %v", err)
		}
	}

	return nil
}

// TODO : need these endpoints
// r.POST("/images", postImageHandler(&database))
// r.DELETE("/images", deleteImageHandler(&database))
func postImageHandler(app_settings common.AppSettings) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10*1000000)
		form, err := c.MultipartForm()
		if err != nil {
			log.Error().Msgf("could not create multipart form: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("request type must be `multipart-form`", err))
			return
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
		file_content_type := file.Header.Get("content-type")
		_, ok := allowed_content_types[file_content_type]
		if !ok {
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

		ext := filepath.Ext(file.Filename)
		// check ext is supported
		_, ok = allowed_extensions[ext]
		if ext == "" || !ok {
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

		// Save lower dimensions of the image if needed
		log.Info().Msgf("creating minified images for %s", image_path)
		err = createMinifiedImages(image_path)
		if err != nil {
			log.Error().Msgf("could not create minified images: %v", err)
		}

		// End saving to filesystem
		c.JSON(http.StatusOK, ImageIdResponse{
			Id: uuid.String(),
		})
	}
}

func deleteImageHandler(app_settings common.AppSettings) func(*gin.Context) {
	return func(c *gin.Context) {
		var delete_image_binding DeleteImageBinding
		err := c.ShouldBindUri(&delete_image_binding)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorRes("no id provided to delete image", err))
			return
		}

		image_path := filepath.Join(app_settings.ImageDirectory, delete_image_binding.Name)
		err = os.Remove(image_path)
		if err != nil {
			log.Warn().Msgf("could not delete stored image file: %v", err)
			// No return because we have to remove the database entry nonetheless.
		}

		c.JSON(http.StatusOK, ImageIdResponse{
			delete_image_binding.Name,
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
