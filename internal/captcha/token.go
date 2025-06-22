package captcha

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type PayloadType string

const (
	AnswerTokenType PayloadType = "answer"
	PassTokenType   PayloadType = "captcha_pass"
)

var DefaultTokenTTL = 5 * time.Minute

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrExpiredToken      = errors.New("expired token")
	ErrInvalidPayload    = errors.New("invalid payload")
	ErrSignatureMismatch = errors.New("signature mismatch")
	ErrWrongTokenType    = errors.New("wrong token type")
)

type TokenManager struct {
	Secret []byte
	TTL    time.Duration
}

func NewTokenManager(secret []byte, ttl time.Duration) *TokenManager {
	return &TokenManager{
		Secret: secret,
		TTL:    ttl,
	}
}

type AnswerPayload struct {
	Type   PayloadType `json:"type"`
	ID     string      `json:"id"`
	Answer string      `json:"answer"`
	Exp    int64       `json:"exp"`
}

type PassPayload struct {
	Type PayloadType `json:"type"`
	Exp  int64       `json:"exp"`
}

func (tm *TokenManager) SignAnswerToken(q *Question) (string, error) {
	payload := AnswerPayload{
		Type:   AnswerTokenType,
		ID:     q.ID,
		Answer: q.Answer,
		Exp:    time.Now().Add(tm.TTL).Unix(),
	}
	return tm.sign(payload)
}

func (tm *TokenManager) SignPassToken() (string, error) {
	payload := PassPayload{
		Type: PassTokenType,
		Exp:  time.Now().Add(tm.TTL).Unix(),
	}
	return tm.sign(payload)
}

func (tm *TokenManager) sign(payload any) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, tm.Secret)
	mac.Write(data)
	sig := mac.Sum(nil)

	token := append(data, sig...)
	return base64.URLEncoding.EncodeToString(token), nil
}

func (tm *TokenManager) VerifyAnswerToken(token string, userAnswer string) (bool, error) {
	raw, err := base64.URLEncoding.DecodeString(token)
	if err != nil || len(raw) <= sha256.Size {
		return false, ErrInvalidToken
	}

	data := raw[:len(raw)-sha256.Size]
	sig := raw[len(raw)-sha256.Size:]

	mac := hmac.New(sha256.New, tm.Secret)
	mac.Write(data)
	if !hmac.Equal(mac.Sum(nil), sig) {
		return false, ErrSignatureMismatch
	}

	var payload AnswerPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return false, ErrInvalidPayload
	}
	if payload.Type != AnswerTokenType {
		return false, ErrWrongTokenType
	}
	if time.Now().Unix() > payload.Exp {
		return false, ErrExpiredToken
	}

	return strings.EqualFold(payload.Answer, userAnswer), nil
}

func (tm *TokenManager) VerifyPassToken(token string) (bool, error) {
	raw, err := base64.URLEncoding.DecodeString(token)
	if err != nil || len(raw) <= sha256.Size {
		return false, ErrInvalidToken
	}

	data := raw[:len(raw)-sha256.Size]
	sig := raw[len(raw)-sha256.Size:]

	mac := hmac.New(sha256.New, tm.Secret)
	mac.Write(data)
	if !hmac.Equal(mac.Sum(nil), sig) {
		return false, ErrSignatureMismatch
	}

	var payload PassPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return false, ErrInvalidPayload
	}
	if payload.Type != PassTokenType {
		return false, ErrWrongTokenType
	}
	if time.Now().Unix() > payload.Exp {
		return false, ErrExpiredToken
	}

	return true, nil
}
