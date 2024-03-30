package app

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/matheusgomes28/urchin/views"
	"github.com/rs/zerolog/log"
)

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func postHandler(c *gin.Context, app_settings common.AppSettings, database database.Database) ([]byte, error) {

	var post_binding common.PostIdBinding

	err := c.ShouldBindUri(&post_binding)

	if err != nil || post_binding.Id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})

		return nil, err
	}

	// Get the post with the ID
	post, err := database.GetPost(post_binding.Id)

	if err != nil || post.Content == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post Not Found"})
		return nil, err
	}

	// Generate HTML page
	post.Content = string(mdToHTML([]byte(post.Content)))
	post_view := views.MakePostPage(post.Title, post.Content)
	html_buffer := bytes.NewBuffer(nil)
	if err = post_view.Render(c, html_buffer); err != nil {
		log.Error().Msgf("could not render: %v", err)
	}

	return html_buffer.Bytes(), nil
}
