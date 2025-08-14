package admin_app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/rs/zerolog/log"
)

// @Summary      Add a new page
// @Description  Adds a new page to the database.
// @Tags         pages
// @Accept       json
// @Produce      json
// @Param        page body AddPageRequest true "Page to add"
// @Success      200 {object} PageResponse
// @Failure      400 {object} common.ErrorResponse "Invalid request body or data"
// @Router       /pages [post]
func postPageHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var add_page_request AddPageRequest
		if c.Request.Body == nil {
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("no request body provided"))
			return
		}

		decoder := json.NewDecoder(c.Request.Body)
		err := decoder.Decode(&add_page_request)
		if err != nil {
			log.Warn().Msgf("invalid page request: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("invalid request body", err))
			return
		}

		err = checkRequiredPageData(add_page_request)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("invalid page data"))
			return
		}

		err = checkRequiredPageData(add_page_request)
		if err != nil {
			log.Error().Msgf("failed to add post required data is missing: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("missing required data", err))
			return
		}

		id, err := database.AddPage(
			add_page_request.Title,
			add_page_request.Content,
			add_page_request.Link,
		)
		if err != nil {
			log.Error().Msgf("failed to add post: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not add post", err))
			return
		}

		c.JSON(http.StatusOK, PageResponse{
			Id:   id,
			Link: add_page_request.Link,
		})
	}
}

// func putPostHandler(database database.Database) func(*gin.Context) {
// 	return func(c *gin.Context) {
// 		var change_post_request ChangePostRequest
// 		decoder := json.NewDecoder(c.Request.Body)
// 		decoder.DisallowUnknownFields()

// 		err := decoder.Decode(&change_post_request)
// 		if err != nil {
// 			log.Warn().Msgf("could not get post from DB: %v", err)
// 			c.JSON(http.StatusBadRequest, common.ErrorRes("invalid request body", err))
// 			return
// 		}

// 		err = database.ChangePost(
// 			change_post_request.Id,
// 			change_post_request.Title,
// 			change_post_request.Excerpt,
// 			change_post_request.Content,
// 		)
// 		if err != nil {
// 			log.Error().Msgf("failed to change post: %v", err)
// 			c.JSON(http.StatusBadRequest, common.ErrorRes("could not change post", err))
// 			return
// 		}

// 		c.JSON(http.StatusOK, PostIdResponse{
// 			change_post_request.Id,
// 		})
// 	}
// }

// func deletePostHandler(database database.Database) func(*gin.Context) {
// 	return func(c *gin.Context) {
// 		var delete_post_binding DeletePostBinding
// 		err := c.ShouldBindUri(&delete_post_binding)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, common.ErrorRes("no id provided to delete post", err))
// 			return
// 		}

// 		rows_affected, err := database.DeletePost(delete_post_binding.Id)
// 		if err != nil {
// 			log.Error().Msgf("failed to delete post: %v", err)
// 			c.JSON(http.StatusBadRequest, common.ErrorRes("could not delete post", err))
// 			return
// 		}

// 		if rows_affected == 0 {
// 			log.Error().Msgf("no post found with id `%d`", delete_post_binding.Id)
// 			c.JSON(http.StatusNotFound, common.MsgErrorRes("no post found"))
// 			return
// 		}

// 		c.JSON(http.StatusOK, PostIdResponse{
// 			delete_post_binding.Id,
// 		})
// 	}
// }

func checkRequiredPageData(add_page_request AddPageRequest) error {
	if strings.TrimSpace(add_page_request.Title) == "" {
		return fmt.Errorf("missing required data 'Title'")
	}

	if strings.TrimSpace(add_page_request.Content) == "" {
		return fmt.Errorf("missing required data 'Content'")
	}

	err := validateLink(add_page_request.Link)
	if err != nil {
		return err
	}

	return nil
}

func validateLink(link string) error {
	for _, char := range link {
		char_val := int(char)
		is_uppercase := (char_val >= 65) && (char_val <= 90)
		is_lowercase := (char_val >= 97) && (char_val <= 122)
		is_sign := (char == '-') || (char == '_')

		if !(is_uppercase || is_lowercase || is_sign) {
			// TODO : what is this conversion?!
			return fmt.Errorf("invalid character in link %s", string(rune(char)))
		}
	}

	return nil
}
