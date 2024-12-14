package app

import (
	"bytes"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// / This function will render the templ component into
// / a gin context's Response Writer
func render(c *gin.Context, status int, template templ.Component) error {
	c.Status(status)
	return template.Render(c.Request.Context(), c.Writer)
}

func renderHtml(c *gin.Context, template templ.Component) ([]byte, error) {
	html_buffer := bytes.NewBuffer(nil)

	err := template.Render(c, html_buffer)
	if err != nil {
		log.Error().Msgf("Could not render index: %v", err)
		return []byte{}, err
	}

	return html_buffer.Bytes(), nil
}
