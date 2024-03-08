package app

import (
	"bytes"
	"fmt"
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

func serveErrorPage(c *gin.Context, error_msg string, error_code int, app_settings common.AppSettings) error {
	// Generate HTML page
	log.Error().Msgf("running the serveErrorPage")
	error_page := views.MakeErrorPage(error_msg, app_settings.AppNavbar.Links)

	// TODO : is there a better function to serve the bytes?
	if err := render(c, error_code, error_page); err == nil {
		log.Error().Msgf("could not render: %v", err)
	}

	return nil
}

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

		err = serveErrorPage(c, "requested invalid post ID", http.StatusBadRequest, app_settings)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		}

		return nil, fmt.Errorf("invalid post id")
	}

	// Get the post with the ID
	post, err := database.GetPost(post_binding.Id)

	if err != nil || post.Content == "" {
		err = serveErrorPage(c, "given post not found", http.StatusNotFound, app_settings)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post Not Found"})
		}
		return nil, err
	}

	// Generate HTML page
	post.Content = string(mdToHTML([]byte(post.Content)))
	post_view := views.MakePostPage(post.Title, post.Content, app_settings.AppNavbar.Links)
	html_buffer := bytes.NewBuffer(nil)
	if err = post_view.Render(c, html_buffer); err != nil {
		log.Error().Msgf("could not render: %v", err)
	}

	return html_buffer.Bytes(), nil
}
