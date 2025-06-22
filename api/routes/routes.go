package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/tracey-io/tracey/api/handlers"
)

func SetupRoutes(app *echo.Echo, handlers *handlers.Handlers) {
	api := app.Group("/v1")
	SetupCaptchaRoutes(api, handlers.CaptchaHandler)
}

func SetupCaptchaRoutes(app *echo.Group, handler *handlers.CaptchaHandler) {
	captcha := app.Group("/captcha")
	captcha.GET("", handler.GetCaptcha)
	captcha.POST("/verify", handler.VerifyCaptcha)
	captcha.POST("/validate", handler.ValidateCaptcha)
}
