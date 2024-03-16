package app

import (
	"bytes"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/matheusgomes28/urchin/views"
	"github.com/rs/zerolog/log"
)

const CACHE_TIMEOUT = 20 * time.Second

type Generator = func(*gin.Context, common.AppSettings, database.Database) ([]byte, error)

func SetupRoutes(app_settings common.AppSettings, database database.Database) *gin.Engine {
	r := gin.Default()
	r.MaxMultipartMemory = 1

	// All cache-able endpoints
	cache := MakeCache(4, time.Minute*10, &TimeValidator{})
	addCachableHandler(r, "GET", "/", homeHandler, &cache, app_settings, database)
	addCachableHandler(r, "GET", "/contact", contactHandler, &cache, app_settings, database)
	addCachableHandler(r, "GET", "/post/:id", postHandler, &cache, app_settings, database)

	// Add the pagination route as a cacheable endpoint
	addCachableHandler(r, "GET", "/page/:num", homeHandler, &cache, app_settings, database)

	// DO not cache as it needs to handlenew form values
	r.POST("/contact-send", makeContactFormHandler())

	r.Static("/static", "./static")
	return r
}

func addCachableHandler(e *gin.Engine, method string, endpoint string, generator Generator, cache *Cache, app_settings common.AppSettings, db database.Database) {

	handler := func(c *gin.Context) {
		// if the endpoint is cached
		cached_endpoint, err := (*cache).Get(c.Request.RequestURI)
		if err == nil {
			c.Data(http.StatusOK, "text/html; charset=utf-8", cached_endpoint.Contents)
			return
		}

		// Before handler call (retrieve from cache)
		html_buffer, err := generator(c, app_settings, db)
		if err != nil {
			log.Error().Msgf("could not generate html: %v", err)
			// TODO : Need a proper error page
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not render HTML",
				"msg":   err.Error(),
			})
			return
		}

		// After handler  (add to cache)
		err = (*cache).Store(c.Request.RequestURI, html_buffer)
		if err != nil {
			log.Warn().Msgf("could not add page to cache: %v", err)
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", html_buffer)
	}

	// Hacky
	if method == "GET" {
		e.GET(endpoint, handler)
	}
	if method == "POST" {
		e.POST(endpoint, handler)
	}
	if method == "DELETE" {
		e.DELETE(endpoint, handler)
	}
	if method == "PUT" {
		e.PUT(endpoint, handler)
	}
}

// / This function will act as the handler for
// / the home page
func homeHandler(c *gin.Context, settings common.AppSettings, db database.Database) ([]byte, error) {
	pageNum := 0 // Default to page 0
	if pageNumQuery := c.Param("num"); pageNumQuery != "" {
		num, err := strconv.Atoi(pageNumQuery)
		if err == nil && num > 0 {
			pageNum = num
		} else {
			log.Error().Msgf("Invalid page number: %s", pageNumQuery)
		}
	}
	limit := 10 // or whatever limit you want
	offset := max((pageNum-1)*limit, 0)

	posts, err := db.GetPosts(limit, offset)
	if err != nil {
		return nil, err
	}

	// if not cached, create the cache
	index_view := views.MakeIndex(posts)
	html_buffer := bytes.NewBuffer(nil)

	err = index_view.Render(c, html_buffer)
	if err != nil {
		log.Error().Msgf("Could not render index: %v", err)
		return []byte{}, err
	}

	return html_buffer.Bytes(), nil
}
