package admin_app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/database"
	"github.com/rs/zerolog/log"
	lua "github.com/yuin/gopher-lua"
)

func getPostHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		// localhost:8080/post/{id}
		var post_binding PostBinding
		if err := c.ShouldBindUri(&post_binding); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "could not get post id",
				"msg":   err.Error(),
			})
			return
		}

		post_id, err := strconv.Atoi(post_binding.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid post id type",
				"msg":   err.Error(),
			})
			return
		}

		post, err := database.GetPost(post_id)
		if err != nil {
			log.Warn().Msgf("could not get post from DB: %v", err)
			c.JSON(http.StatusNotFound, gin.H{
				"error": "post id not found",
				"msg":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":      post.Id,
			"title":   post.Title,
			"excerpt": post.Excerpt,
			"content": post.Content,
		})
	}
}

func postPostHandler(database database.Database, shortcode_handlers map[string]*lua.LState) func(*gin.Context) {
	return func(c *gin.Context) {
		var add_post_request AddPostRequest
		decoder := json.NewDecoder(c.Request.Body)
		err := decoder.Decode(&add_post_request)

		if err != nil {
			log.Warn().Msgf("could not get post from DB: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
				"msg":   err.Error(),
			})
			return
		}

		transformed_content, err := transformContent(add_post_request.Content, shortcode_handlers)
		if err != nil {
			log.Warn().Msgf("could not transform post: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
				"msg":   err.Error(),
			})
			return
		}

		id, err := database.AddPost(
			add_post_request.Title,
			add_post_request.Excerpt,
			transformed_content,
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
			"id": id,
		})
	}
}

func putPostHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var change_post_request ChangePostRequest
		decoder := json.NewDecoder(c.Request.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&change_post_request)
		if err != nil {
			log.Warn().Msgf("could not get post from DB: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
				"msg":   err.Error(),
			})
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
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "could not change post",
				"msg":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": change_post_request.Id,
		})
	}
}

func deletePostHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var delete_post_request DeletePostRequest
		decoder := json.NewDecoder(c.Request.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&delete_post_request)
		if err != nil {
			log.Warn().Msgf("could not delete post: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
				"msg":   err.Error(),
			})
			return
		}

		err = database.DeletePost(delete_post_request.Id)
		if err != nil {
			log.Error().Msgf("failed to delete post: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "could not delete post",
				"msg":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": delete_post_request.Id,
		})
	}
}

// partitionString will partition the strings by
// removing the given ranges
func partitionString(text string, indexes [][]int) []string {

	partitions := make([]string, 0)
	start := 0
	for _, window := range indexes {
		partitions = append(partitions, text[start:window[0]])
		start = window[1]
	}

	partitions = append(partitions, text[start:len(text)-1])
	return partitions
}

func shortcodeToMarkdown(shortcode string, shortcode_handlers map[string]*lua.LState) (string, error) {
	key_value := strings.Split(shortcode, ":")

	key := key_value[0]
	values := key_value[1:]

	if handler, found := shortcode_handlers[key]; found {

		// Need to quote all values for a valid lua syntax
		quoted_values := make([]string, 0)
		for _, value := range values {
			quoted_values = append(quoted_values, fmt.Sprintf("%q", value))
		}

		err := handler.DoString(fmt.Sprintf(`result = HandleShortcode({%s})`, strings.Join(quoted_values, ",")))
		if err != nil {
			return "", fmt.Errorf("error running %s shortcode: %v", key, err)
		}

		value := handler.GetGlobal("result")
		if ret_type := value.Type().String(); ret_type != "string" {
			return "", fmt.Errorf("error running %s shortcode: invalid return type %s", key, ret_type)
		} else if ret_type == "" {
			return "", fmt.Errorf("error running %s shortcode: returned empty string", key)
		}

		return value.String(), nil
	}

	return "", fmt.Errorf("unsupported shortcode: %s", key)
}

func transformContent(content string, shortcode_handlers map[string]*lua.LState) (string, error) {
	// Find all the occurences of {{ and }}
	regex, _ := regexp.Compile(`{{[\w.-]+(:[\w.-]+)+}}`)

	shortcodes := regex.FindAllStringIndex(content, -1)
	partitions := partitionString(content, shortcodes)

	builder := strings.Builder{}
	i := 0
	for i, shortcode := range shortcodes {
		builder.WriteString(partitions[i])

		markdown, err := shortcodeToMarkdown(content[shortcode[0]+2:shortcode[1]-2], shortcode_handlers)
		if err != nil {
			log.Error().Msgf("%v", err)
			markdown = ""
		}
		builder.WriteString(markdown)
	}

	// Guaranteed to have +1 than the number of
	// shortcodes by algorithm
	builder.WriteString(partitions[i+1])

	return builder.String(), nil
}
