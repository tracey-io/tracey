package captcha

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

var (
	ErrInvalidPOW = errors.New("invalid proof of work")
	ErrExpiredPOW = errors.New("expired proof of work")
)

var DefaultPOWConfig = &POWConfig{
	Difficulty: 20,
	TTL:        2 * time.Minute,
}

type POWConfig struct {
	Difficulty int
	TTL        time.Duration
}

type POWChallenge struct {
	Nonce      string `json:"nonce"`
	Difficulty int    `json:"difficulty"`
	Timestamp  int64  `json:"timestamp"`
}

type POWManager struct {
	PowConfig *POWConfig
}

func NewPOWManager(powConfig *POWConfig) *POWManager {
	return &POWManager{
		PowConfig: powConfig,
	}
}

func (p *POWManager) GenerateChallenge() *POWChallenge {
	return &POWChallenge{
		Nonce:      uuid.NewString(),
		Difficulty: p.PowConfig.Difficulty,
		Timestamp:  time.Now().Unix(),
	}
}

func (p *POWManager) Verify(nonce string, counter int, issuedAt int64) error {
	if time.Since(time.Unix(issuedAt, 0)) > p.PowConfig.TTL {
		return ErrExpiredPOW
	}

	input := fmt.Sprintf("%s:%d", nonce, counter)
	hash := sha256.Sum256([]byte(input))

	if !hasLeadingZeroBits(hash[:], p.PowConfig.Difficulty) {
		return ErrInvalidPOW
	}

	return nil
}

func hasLeadingZeroBits(hash []byte, bits int) bool {
	fullBytes := bits / 8
	remainingBits := bits % 8

	for i := 0; i < fullBytes; i++ {
		if hash[i] != 0x00 {
			return false
		}
	}

	if remainingBits > 0 {
		mask := byte(0xFF << (8 - remainingBits))
		if hash[fullBytes]&mask != 0x00 {
			return false
		}
	}

	return true
}
