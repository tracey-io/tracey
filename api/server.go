package api

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tracey-io/tracey/api/handlers"
	"github.com/tracey-io/tracey/api/routes"
	"github.com/tracey-io/tracey/internal/captcha"
	"golang.org/x/crypto/acme/autocert"
	"os"
	"os/signal"
	"time"
)

func StartServer(serverConfig *ServerConfig) error {
	app := echo.New()

	app.HideBanner = true
	app.HidePort = true

	app.Use(middleware.Recover())
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
		Output: app.Logger.Output(),
	}))

	questionManager := captcha.NewQuestionManager(captcha.DefaultCategories)
	tokenManager := captcha.NewTokenManager([]byte(serverConfig.SecretKey), captcha.DefaultTokenTTL)
	powManager := captcha.NewPOWManager(captcha.DefaultPOWDifficulty, captcha.DefaultPOWTTL)
	captchaService := captcha.NewService(questionManager, tokenManager, powManager)
	captchaHandler := handlers.NewCaptchaHandler(captchaService)

	serverHandlers := handlers.NewHandlers().
		AddCaptchaHandler(captchaHandler)

	routes.SetupRoutes(app, serverHandlers)

	serverAddress := serverConfig.Address.String()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	defer stop()

	go func() {
		log.Debug("Starting server", "host", serverConfig.Address.Host, "port", serverConfig.Address.Port, "environment", serverConfig.Environment)

		var err error

		switch serverConfig.Environment {
		case EnvironmentDev:
			app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
				AllowOrigins: []string{
					"http://localhost:3000",
					"http://127.0.0.1:3000",
				},
				AllowMethods: []string{
					echo.GET,
					echo.POST,
					echo.PUT,
					echo.DELETE,
					echo.OPTIONS,
				},
				AllowHeaders: []string{
					echo.HeaderOrigin,
					echo.HeaderContentType,
					echo.HeaderAccept,
					echo.HeaderAuthorization,
				},
				AllowCredentials: true,
			}))

			err = app.Start(serverAddress)
		case EnvironmentProd:
			app.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
			err = app.StartAutoTLS(serverAddress)
		default:
			err = ErrEnvironmentNotSupported
		}

		if err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	return nil
}
