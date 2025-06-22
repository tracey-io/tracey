package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/tracey-io/tracey/internal/captcha"
	"net/http"
)

type CaptchaHandler struct {
	CaptchaService *captcha.Service
}

func NewCaptchaHandler(captchaService *captcha.Service) *CaptchaHandler {
	return &CaptchaHandler{
		CaptchaService: captchaService,
	}
}

func (h *CaptchaHandler) GetCaptcha(ctx echo.Context) error {
	challenge, err := h.CaptchaService.Generate()

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"success": false,
			"message": "failed to generate captcha",
		})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"success": true,
		"data": echo.Map{
			"id":      challenge.ID,
			"prompt":  challenge.Prompt,
			"options": challenge.Options,
			"token":   challenge.Token,
			"pow": echo.Map{
				"nonce":      challenge.POWNonce,
				"difficulty": challenge.Difficulty,
				"timestamp":  challenge.Timestamp,
			},
		},
	})
}

func (h *CaptchaHandler) VerifyCaptcha(ctx echo.Context) error {
	var body struct {
		Answer    string `json:"answer"`
		Token     string `json:"token"`
		Nonce     string `json:"nonce"`
		Counter   int    `json:"counter"`
		Timestamp int64  `json:"timestamp"`
	}

	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": "failed parsing request body",
		})
	}

	if body.Answer == "" || body.Token == "" || body.Nonce == "" || body.Counter == 0 || body.Timestamp == 0 {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	passToken, err := h.CaptchaService.Verify(body.Answer, body.Token, body.Nonce, body.Counter, body.Timestamp)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": "failed to verify captcha",
		})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"success": true,
		"data": echo.Map{
			"token": passToken,
		},
	})
}

func (h *CaptchaHandler) ValidateCaptcha(ctx echo.Context) error {
	var body struct {
		PassToken string `json:"passToken"`
	}

	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": "failed parsing request body",
		})
	}

	passToken := body.PassToken

	if passToken == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	valid, err := h.CaptchaService.Validate(body.PassToken)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"success": false,
			"message": "failed verifying token",
		})
	}

	if !valid {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": "invalid captcha token",
		})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"success": true,
	})
}
