package app

import (
	"bytes"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/matheusgomes28/urchin/views"
)

type PostBinding struct {
	Id string `uri:"id" binding:"required"`
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

func postHandler(c *gin.Context, app_settings common.AppSettings, database *database.Database) ([]byte, error) {
	// localhost:8080/post/{id}

	var post_binding PostBinding
	if err := c.ShouldBindUri(&post_binding); err != nil {
		return nil, err
	}
	
	// Get the post with the ID 
	post_id, err := strconv.Atoi(post_binding.Id)
	if err != nil {
		return nil, err
	}

	post, err := database.GetPost(post_id)
	if err != nil {
		return nil, err
	}

	// Generate HTML page
	post.Content = string(mdToHTML([]byte(post.Content)))
	post_view := views.MakePostPage(post.Title, post.Content)
	html_buffer := bytes.NewBuffer(nil)
	post_view.Render(c, html_buffer)

	return html_buffer.Bytes(), nil
}