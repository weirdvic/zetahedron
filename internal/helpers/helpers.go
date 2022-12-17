package helpers

import (
	"errors"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

const (
	alphabet      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	largestuint64 = 18446744073709551615
)

type (
	Request struct {
		URL    string        `json:"url" validate:"required,url"`
		Slug   string        `json:"url_slug" validate:"printascii,max=32"`
		Expiry time.Duration `json:"expiry"`
	}

	Response struct {
		URL    string        `json:"url"`
		Slug   string        `json:"url_slug"`
		Expiry time.Duration `json:"expiry"`
	}

	CustomValidator struct {
		Validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// Base62Encode
func Base62Encode(number uint64) string {
	length := len(alphabet)
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(10)
	for ; number > 0; number = number / uint64(length) {
		encodedBuilder.WriteByte(alphabet[(number % uint64(length))])
	}

	return encodedBuilder.String()
}

// Base62Decode
func Base62Decode(encodedString string) (uint64, error) {
	var number uint64
	length := len(alphabet)
	for i, symbol := range encodedString {
		alphabeticPosition := strings.IndexRune(alphabet, symbol)
		if alphabeticPosition == -1 {
			return uint64(alphabeticPosition), errors.New("cannot find symbol in alphabet")
		}
		number += uint64(alphabeticPosition) * uint64(math.Pow(float64(length), float64(i)))
	}
	return number, nil
}

// EnforceHTTPS checks if URL starts with https://
// and replaces/adds schema if needed
func EnforceHTTPS(url string) string {
	if url[:5] == "http:" {
		return strings.Replace(url, "http:", "https:", 1)
	} else if url[:6] != "https:" {
		return "https://" + url
	}
	return url
}

// IsOurDomain checks if provided URL is different
// from our own domain name
func IsOurDomain(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return true
	}
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]
	if newURL != os.Getenv("DOMAIN") {
		return false
	}
	return false
}
