package main

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/weirdvic/zetahedron/internal/database"
	h "github.com/weirdvic/zetahedron/internal/helpers"
)

// app.echo.POST("/shorten", shortenURL)
// endpoint to shorten the URL and put it to the DB
func (app *application) shortenURL(c echo.Context) (err error) {
	// Parse and validate the request
	body := new(h.Request)
	if err = c.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(body); err != nil {
		return err
	}

	// Do not short URLs leading to our domain
	if !h.IsOurDomain(body.URL) {
		return c.JSON(http.StatusBadRequest, &map[string]string{
			"error": "Can't do that :)",
		})
	}

	// Switch to HTTPS if needed
	body.URL = h.EnforceHTTPS(body.URL)

	// Make random key for storing in the DB or use body.Slug
	var urlKey string
	if body.Slug == "" {
		urlKey = h.Base62Encode(rand.Uint64())
	} else {
		urlKey = body.Slug
	}

	// Checking if url_slug is already used in DB
	r := database.CreateClient(0)
	defer r.Close()
	val, _ := r.Get(database.Ctx, urlKey).Result()

	if val != "" {
		return c.JSON(http.StatusBadRequest, &map[string]string{
			"error": "URL slug is already in use",
		})
	}
	// Set default expiry time to 24 hours
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = r.Set(database.Ctx, urlKey, body.URL, body.Expiry*3600*time.Second).Err()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &map[string]string{"error": "Unable to connect to DB server"})
	}

	resp := h.Response{
		URL:    body.URL,
		Slug:   "",
		Expiry: body.Expiry,
	}

	resp.Slug = os.Getenv("DOMAIN") + "/" + urlKey
	return c.JSON(http.StatusOK, resp)
}

// app.echo.GET("/:slug", resolveURL)
func (app *application) resolveURL(c echo.Context) (err error) {
	slug := c.Param("slug")
	r := database.CreateClient(0)
	defer r.Close()

	url, err := r.Get(database.Ctx, slug).Result()
	if err == redis.Nil {
		return c.JSON(http.StatusNotFound, &map[string]string{"error": "URL slug not found in DB"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, &map[string]string{"error": "Internal server error"})
	}
	return c.Redirect(http.StatusMovedPermanently, url)
}
