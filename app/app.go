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
	cache := makeCache(4, time.Minute*10)
	addCachableHandler(r, "GET", "/", homeHandler, &cache, app_settings, database)
	addCachableHandler(r, "GET", "/contact", contactHandler, &cache, app_settings, database)
	addCachableHandler(r, "GET", "/post/:id", postHandler, &cache, app_settings, database)

	// DO not cache as it needs to handlenew form values
	r.POST("/contact-send", makeContactFormHandler())

	// this will setup the route for pagination
	r.GET("/page/:num", func(c *gin.Context) {
		homeHandler(c, app_settings, database)
	})

	r.Static("/static", "./static")
	return r
}

func addCachableHandler(e *gin.Engine, method string, endpoint string, generator Generator, cache *Cache, app_settings common.AppSettings, db database.Database) {

	handler := func(c *gin.Context) {
		// if the endpoint is cached
		cached_endpoint, err := (*cache).Get(c.Request.RequestURI)
		if err == nil {
			c.Data(http.StatusOK, "text/html; charset=utf-8", cached_endpoint.contents)
			return
		}

		// Before handler call (retrieve from cache)
		html_buffer, err := generator(c, app_settings, db)
		if err != nil {
			log.Error().Msgf("could not generate html: %v", err)
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
	pageNumQuery := c.Param("num")
	pageNum, err := strconv.Atoi(pageNumQuery)
	if err != nil {
		return nil, err
	}

	limit := 10 // or whatever limit you want
	offset := (pageNum - 1) * limit

	posts, err := db.GetPosts(limit, offset)
	if err != nil {
		return nil, err
	}

	// if not cached, create the cache
	index_view := views.MakeIndex(posts)
	html_buffer := bytes.NewBuffer(nil)

	err = index_view.Render(c, html_buffer)
	if err != nil {
		log.Error().Msgf("could not render index: %v", err)
		return []byte{}, err
	}

	return html_buffer.Bytes(), nil
}
