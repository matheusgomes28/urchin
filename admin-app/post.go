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

// @Summary      Get a post by ID
// @Description  Retrieves a post based on its ID.
// @Tags         posts
// @Produce      json
// @Param        id path int true "Post ID"
// @Success      200 {object} GetPostResponse
// @Failure      400 {object} common.ErrorResponse "Invalid ID or post not found"
// @Router       /posts/{id} [get]
func getPostHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		// localhost:8080/post/{id}
		var post_binding common.PostIdBinding
		if err := c.ShouldBindUri(&post_binding); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not get post id", err))
			return
		}

		post, err := database.GetPost(post_binding.Id)
		if err != nil {
			log.Warn().Msgf("could not get post from DB: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("post id not found", err))
			return
		}

		c.JSON(http.StatusOK, GetPostResponse{
			post.Id,
			post.Title,
			post.Excerpt,
			post.Content,
		})
	}
}

// @Summary      Add a new post
// @Description  Adds a new post to the database.
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        post body AddPostRequest true "Post to add"
// @Success      200 {object} PostIdResponse
// @Failure      400 {object} common.ErrorResponse "Invalid request body or missing data"
// @Router       /posts [post]
func postPostHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var add_post_request AddPostRequest
		if c.Request.Body == nil {
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("no request body provided"))
			return
		}
		decoder := json.NewDecoder(c.Request.Body)
		err := decoder.Decode(&add_post_request)

		if err != nil {
			log.Warn().Msgf("invalid post request: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("invalid request body", err))
			return
		}

		err = checkRequiredData(add_post_request)
		if err != nil {
			log.Error().Msgf("failed to add post required data is missing: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("missing required data", err))
			return
		}

		id, err := database.AddPost(
			add_post_request.Title,
			add_post_request.Excerpt,
			add_post_request.Content,
		)
		if err != nil {
			log.Error().Msgf("failed to add post: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not add post", err))
			return
		}

		c.JSON(http.StatusOK, PostIdResponse{
			id,
		})
	}
}

// @Summary      Update an existing post
// @Description  Updates an existing post with new data.
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        post body ChangePostRequest true "Post data to update"
// @Success      200 {object} PostIdResponse
// @Failure      400 {object} common.ErrorResponse "Invalid request body or could not change post"
// @Router       /posts [put]
func putPostHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var change_post_request ChangePostRequest
		decoder := json.NewDecoder(c.Request.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&change_post_request)
		if err != nil {
			log.Warn().Msgf("could not get post from DB: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("invalid request body", err))
			return
		}

		err = database.ChangePost(
			change_post_request.Id,
			change_post_request.Title,
			change_post_request.Excerpt,
			change_post_request.Content,
		)
		if err != nil {
			log.Error().Msgf("failed to change post: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not change post", err))
			return
		}

		c.JSON(http.StatusOK, PostIdResponse{
			change_post_request.Id,
		})
	}
}

// @Summary      Delete a post
// @Description  Deletes a post by its ID.
// @Tags         posts
// @Produce      json
// @Param        id path int true "Post ID"
// @Success      200 {object} PostIdResponse
// @Failure      400 {object} common.ErrorResponse "Invalid ID provided"
// @Failure      404 {object} common.ErrorResponse "Post not found"
// @Router       /posts/{id} [delete]
func deletePostHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var delete_post_binding DeletePostBinding
		err := c.ShouldBindUri(&delete_post_binding)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorRes("no id provided to delete post", err))
			return
		}

		rows_affected, err := database.DeletePost(delete_post_binding.Id)
		if err != nil {
			log.Error().Msgf("failed to delete post: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not delete post", err))
			return
		}

		if rows_affected == 0 {
			log.Error().Msgf("no post found with id `%d`", delete_post_binding.Id)
			c.JSON(http.StatusNotFound, common.MsgErrorRes("no post found"))
			return
		}

		c.JSON(http.StatusOK, PostIdResponse{
			delete_post_binding.Id,
		})
	}
}

func checkRequiredData(addPostRequest AddPostRequest) error {
	if strings.TrimSpace(addPostRequest.Title) == "" {
		return fmt.Errorf("missing required data 'Title'")
	}

	if strings.TrimSpace(addPostRequest.Excerpt) == "" {
		return fmt.Errorf("missing required data 'Excerpt'")
	}

	if strings.TrimSpace(addPostRequest.Content) == "" {
		return fmt.Errorf("missing required data 'Content'")
	}

	return nil
}
