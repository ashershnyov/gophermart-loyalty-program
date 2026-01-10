package jwt

import (
	"fmt"
	"time"

	"github.com/caarlos0/env"
	"github.com/golang-jwt/jwt/v5"
)

type UserIDCtxKey string

const CtxKey UserIDCtxKey = "userID"

const (
	defaultTTL    = 60 * time.Minute
	defaultSecret = "defSec"
)

// UserIDKey is a key used in token claims to store UserID.
const UserIDKey = "user_id"

// Config defines the config of the JWT token generator.
type Config struct {
	TokenTTL    time.Duration `env:"JWT_TTL" envDefault:"0"`
	TokenSecret string        `env:"JWT_SECRET" envDefault:""`
}

// NewConfig creates a new config for JWT token generator.
func NewConfig() (Config, error) {
	c := Config{}
	env.Parse(&c)

	if c.TokenSecret == "" {
		c.TokenSecret = defaultSecret
	}

	if c.TokenTTL == 0 {
		c.TokenTTL = defaultTTL
	}

	return c, nil
}

// Generator is a JWT token provider.
type Generator struct {
	cfg Config
}

// NewGenerator returns a new JWT token generator.
func NewGenerator() (*Generator, error) {
	cfg, err := NewConfig()
	if err != nil {
		return nil, fmt.Errorf("failed creating JWT provider: %w", err)
	}

	return &Generator{
		cfg: cfg,
	}, nil
}

// GetToken returns a token for a requested userID.
func (g *Generator) GetToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		UserIDKey: userID,
		"exp":     time.Now().Add(g.cfg.TokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(g.cfg.TokenSecret))
}

// ParseToken parses token to JWT format from the incoming string.
func (g *Generator) ParseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(g.cfg.TokenSecret), nil
	})
}
