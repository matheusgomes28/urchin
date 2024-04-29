package helpers

import (
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/textproto"
)

func CreateFormImagePart(writer *multipart.Writer, fieldname string, filename string, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldname, filename))
	h.Set("Content-Type", contentType)
	return writer.CreatePart(h)
}

func CreateTextFormHeader(writer *multipart.Writer, fieldname string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"`, fieldname))
	return writer.CreatePart(h)
}

// Creating an image in memory for testing: https://yourbasic.org/golang/create-image/
func CreateImage() image.Image {
	width := 1
	height := 1

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	return image.NewRGBA(image.Rectangle{upLeft, lowRight})
}
