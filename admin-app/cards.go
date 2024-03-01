package admin_app

import (
	"encoding/json"
	"net/http"

	"github.com/fossoreslp/go-uuid-v4"
	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/database"
	"github.com/rs/zerolog/log"
)

type AddCardRequest struct {
	ImageLocation string `json:"image_location"`
	JsonData      string `json:"json_data"`
	SchemaName    string `json:"json_schema"`
}

// func getCardHandler(database database.Database) func(*gin.Context) {
// 	return func(c *gin.Context) {
// 		// localhost:8080/post/{id}
// 		var post_binding PostBinding
// 		if err := c.ShouldBindUri(&post_binding); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": "could not get post id",
// 				"msg":   err.Error(),
// 			})
// 			return
// 		}

// 		post_id, err := strconv.Atoi(post_binding.Id)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": "invalid post id type",
// 				"msg":   err.Error(),
// 			})
// 			return
// 		}

// 		post, err := database.GetPost(post_id)
// 		if err != nil {
// 			log.Warn().Msgf("could not get post from DB: %v", err)
// 			c.JSON(http.StatusNotFound, gin.H{
// 				"error": "post id not found",
// 				"msg":   err.Error(),
// 			})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{
// 			"id":      post.Id,
// 			"title":   post.Title,
// 			"excerpt": post.Excerpt,
// 			"content": post.Content,
// 		})
// 	}
// }

func postCardHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var add_card_request AddCardRequest
		decoder := json.NewDecoder(c.Request.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&add_card_request)

		if err != nil {
			log.Warn().Msgf("invalid post card request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
				"msg":   err.Error(),
			})
			return
		}

		card_uuid, err := uuid.New()
		if err != nil {
			log.Error().Msgf("could not create the UUID: %v", err)
			return
		}

		err = database.AddCard(
			card_uuid.String(),
			add_card_request.ImageLocation,
			add_card_request.JsonData,
			add_card_request.SchemaName,
		)
		if err != nil {
			log.Error().Msgf("failed to add post: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "could not add post",
				"msg":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": card_uuid.String(),
		})
	}
}

// func putCardHandler(database database.Database) func(*gin.Context) {
// 	return func(c *gin.Context) {
// 		var change_post_request ChangePostRequest
// 		decoder := json.NewDecoder(c.Request.Body)
// 		decoder.DisallowUnknownFields()

// 		err := decoder.Decode(&change_post_request)
// 		if err != nil {
// 			log.Warn().Msgf("could not get post from DB: %v", err)
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": "invalid request body",
// 				"msg":   err.Error(),
// 			})
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
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": "could not change post",
// 				"msg":   err.Error(),
// 			})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{
// 			"id": change_post_request.Id,
// 		})
// 	}
// }

// func deleteCardHandler(database database.Database) func(*gin.Context) {
// 	return func(c *gin.Context) {
// 		var delete_post_request DeletePostRequest
// 		decoder := json.NewDecoder(c.Request.Body)
// 		decoder.DisallowUnknownFields()

// 		err := decoder.Decode(&delete_post_request)
// 		if err != nil {
// 			log.Warn().Msgf("could not delete post: %v", err)
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": "invalid request body",
// 				"msg":   err.Error(),
// 			})
// 			return
// 		}

// 		err = database.DeletePost(delete_post_request.Id)
// 		if err != nil {
// 			log.Error().Msgf("failed to delete post: %v", err)
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": "could not delete post",
// 				"msg":   err.Error(),
// 			})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{
// 			"id": delete_post_request.Id,
// 		})
// 	}
// }
