/*
	    There's a lot of common functionality between the image
		and gallery handlers, so this file is meant to share those
		functionalities.
*/

package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog/log"
)

// This list contains the valid file
// extensions for an image.
var ValidImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

func populateImageMetadata(metadata_path string, app_settings AppSettings) (Image, error) {

	// Check if a json metadata file exists
	metadata_contents, err := os.ReadFile(path.Join(app_settings.ImageDirectory, metadata_path))
	if err != nil {
		return Image{}, fmt.Errorf("could not read metadata for image `%s`", metadata_path)
	}

	var image Image
	err = json.Unmarshal(metadata_contents, &image)

	if err != nil {
		return Image{}, fmt.Errorf("could not deserailize metadata for image `%s`", metadata_path)
	}

	ext := path.Ext(image.Filename)

	// Checking for the existence of a value in a map takes O(1) and therefore it's faster than
	// iterating over a string slice
	_, ok := ValidImageExtensions[ext]
	if !ok {
		return Image{}, fmt.Errorf("image type provided in metadata `%s` is not supported: `%s`", metadata_path, image.Filename)
	}

	filepath := path.Join("/images/data", image.Filename)

	metadata_uuid := strings.TrimSuffix(metadata_path, ext)
	image.Ext = ext
	image.Uuid = metadata_uuid
	image.Filepath = filepath
	return image, nil
}

// Given a list of files, this function will return
// a filtered list of valid images, with the page number
// and page size taken as pagination arguments.
//
// paths must be a list of strings referencing the metadata file for an image.
// page_size must be a non-negative number greater than zero.
// page_num must be a non-negative number greater than 0.
func GetImages(paths []string, page_size, page_num int, app_settings AppSettings) ([]Image, error) {

	if page_num <= 0 {
		return []Image{}, fmt.Errorf("invalid `page_num` (%d) given", page_num)
	}

	offset := (page_num - 1) * page_size

	num_paths := len(paths)
	if offset >= num_paths {
		return []Image{}, fmt.Errorf("invalid pagination settings: `page_size` (%d) and `page_num` (%d)", page_size, page_num)
	}

	valid_images := make([]Image, 0)
	for _, metadata_path := range paths {
		image, err := populateImageMetadata(metadata_path, app_settings)
		if err != nil {
			log.Warn().Msgf("skipping image defined in metadata path `%s`: %v", metadata_path, err)
			continue
		}
		valid_images = append(valid_images, image)
	}
	return valid_images, nil
}
