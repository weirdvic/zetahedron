package main

import (
	"flag"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	h "github.com/weirdvic/zetahedron/internal/helpers"
	"go.uber.org/zap"
)

type application struct {
	echo          *echo.Echo
	http_endpoint *string
	logger        *zap.SugaredLogger
}

func (app *application) init() {
	// Set up middleware
	app.echo.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))
	// Set up validator
	app.echo.Validator = &h.CustomValidator{Validator: validator.New()}
	// Set up routing
	app.echo.POST("/shorten", app.shortenURL)
	app.echo.GET("/:url", app.resolveURL)
}

func main() {
	// Create middleware logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	http_endpoint := flag.String("addr", ":1323", "HTTP server address")
	flag.Parse()

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Could not load environment file.")
	}

	// Create application instance
	app := &application{
		echo:          echo.New(),
		http_endpoint: http_endpoint,
		logger:        logger.Sugar(),
	}

	// Attach logger for requests
	app.echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	}))

	// Init server instance
	app.init()

	// Start the server
	app.echo.Logger.Fatal(app.echo.Start(*app.http_endpoint))
}
