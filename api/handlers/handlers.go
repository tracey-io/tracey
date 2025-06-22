package handlers

type Handlers struct {
	CaptchaHandler *CaptchaHandler
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) AddCaptchaHandler(captchaHandler *CaptchaHandler) *Handlers {
	h.CaptchaHandler = captchaHandler
	return h
}
