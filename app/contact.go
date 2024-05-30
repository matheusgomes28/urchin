package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/matheusgomes28/urchin/views"
	"github.com/rs/zerolog/log"
)

// Change this in case Google decides to deprecate
// the reCAPTCHA validation endpoint
const RECAPTCHA_VERIFY_URL string = "https://www.google.com/recaptcha/api/siteverify"

type RecaptchaResponse struct {
	Success   bool    `json:"success"`
	Score     float32 `json:"score"`
	Timestamp string  `json:"challenge_ts"`
	Hostname  string  `json:"hostname"`
}

func verifyRecaptcha(recaptcha_secret string, recaptcha_response string) error {
	// Validate that the recaptcha response was actually
	// not a bot by checking the success rate
	recaptcha_response_data, err := http.PostForm(RECAPTCHA_VERIFY_URL, url.Values{
		"secret":   {recaptcha_secret},
		"response": {recaptcha_response},
	})
	if err != nil {
		err_str := fmt.Sprintf("could not do recaptcha post request: %s", err)
		return fmt.Errorf("%s: %s", err_str, err)
	}
	defer recaptcha_response_data.Body.Close()

	if recaptcha_response_data.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid recaptcha response: %s", recaptcha_response_data.Status)
	}
	var recaptcha_answer RecaptchaResponse
	recaptcha_response_data_buffer, _ := io.ReadAll(recaptcha_response_data.Body)
	err = json.Unmarshal(recaptcha_response_data_buffer, &recaptcha_answer)
	if err != nil {
		return fmt.Errorf("could not parse recaptcha response: %s", err)
	}

	if !recaptcha_answer.Success || (recaptcha_answer.Score < 0.9) {
		return fmt.Errorf("could not validate recaptcha")
	}

	return nil
}

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("could not parse email: %s", email)
	}

	return nil
}

func renderErrorPage(c *gin.Context, email string, err error) error {
	if err = render(c, http.StatusOK, views.MakeContactFailure(email, err.Error())); err != nil {
		log.Error().Msgf("could not render error page: %v", err)
	}
	return err
}

func logError(err error) {
	if err != nil {
		log.Error().Msgf("%v", err)
	}
}

func makeContactFormHandler(app_settings common.AppSettings) func(*gin.Context) {
	return func(c *gin.Context) {
		if err := c.Request.ParseForm(); err != nil {
			log.Error().Msgf("could not parse form %v", err)
			if err = render(c, http.StatusOK, views.MakeContactFailure("unknown", err.Error())); err != nil {
				log.Error().Msgf("could not render %v", err)
			}
			return
		}

		email := c.Request.FormValue("email")
		name := c.Request.FormValue("name")
		message := c.Request.FormValue("message")
		recaptcha_response := c.Request.FormValue("g-recaptcha-response")

		// Make the request to Google's API only if user
		// configured recatpcha settings
		if (len(app_settings.RecaptchaSecret) > 0) && (len(app_settings.RecaptchaSiteKey) > 0) {
			err := verifyRecaptcha(app_settings.RecaptchaSecret, recaptcha_response)
			if err != nil {
				log.Error().Msgf("%v", err)
				defer logError(renderErrorPage(c, email, err))
				return
			}
		}

		err := validateEmail(email)
		if err != nil {
			log.Error().Msgf("%v", err)
			defer logError(renderErrorPage(c, email, err))
			return
		}

		// Make sure name and message is reasonable
		if len(name) > 200 {
			if err = render(c, http.StatusOK, views.MakeContactFailure(email, "name too long (200 chars max)")); err != nil {
				log.Error().Msgf("could not render: %v", err)
				logError(renderErrorPage(c, email, err))
			}
			return
		}

		if len(message) > 10000 {
			if err = render(c, http.StatusOK, views.MakeContactFailure(email, "message too long (1000 chars max)")); err != nil {
				log.Error().Msgf("could not render: %v", err)
				logError(renderErrorPage(c, email, err))
			}
			return
		}

		if err = render(c, http.StatusOK, views.MakeContactSuccess(email, name)); err != nil {
			log.Error().Msgf("could not render: %v", err)
		}
	}
}

// TODO : This is a duplicate of the index handler... abstract
func contactHandler(c *gin.Context, app_settings common.AppSettings, db database.Database) ([]byte, error) {
	index_view := views.MakeContactPage(app_settings.AppNavbar.Links, app_settings.RecaptchaSiteKey)
	html_buffer := bytes.NewBuffer(nil)
	if err := index_view.Render(c, html_buffer); err != nil {
		log.Error().Msgf("could not render: %v", err)
	}

	return html_buffer.Bytes(), nil
}
