package jwt

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newConfigWithValues(secret string, ttl time.Duration) (Config, error) {
	return Config{
		TokenSecret: secret,
		TokenTTL:    ttl,
	}, nil
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name           string
		envSetup       func()
		expectedTTL    time.Duration
		expectedSecret string
	}{
		{
			name:           "default TTL with secret",
			envSetup:       func() { t.Setenv("JWT_SECRET", "testsecret") },
			expectedTTL:    defaultTTL,
			expectedSecret: "testsecret",
		},
		{
			name:           "custom TTL with secret",
			envSetup:       func() { t.Setenv("JWT_SECRET", "testsecret"); t.Setenv("JWT_TTL", "30m") },
			expectedTTL:    30 * time.Minute,
			expectedSecret: "testsecret",
		},
		{
			name:           "no secret",
			envSetup:       func() { os.Unsetenv("JWT_SECRET"); os.Unsetenv("JWT_TTL") },
			expectedTTL:    defaultTTL,
			expectedSecret: defaultSecret,
		},
		{
			name:           "empty secret",
			envSetup:       func() { t.Setenv("JWT_SECRET", ""); os.Unsetenv("JWT_TTL") },
			expectedTTL:    defaultTTL,
			expectedSecret: defaultSecret,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.envSetup()
			cfg, err := NewConfig()
			require.NoError(t, err)
			assert.Equal(t, tt.expectedTTL, cfg.TokenTTL)
			assert.Equal(t, tt.expectedSecret, cfg.TokenSecret)
		})
	}
}

func TestNewGenerator(t *testing.T) {
	gen, _ := NewGenerator()
	assert.NotNil(t, gen)
	assert.Equal(t, Config{
		TokenTTL:    defaultTTL,
		TokenSecret: defaultSecret,
	}, gen.cfg)
}

func TestGenerator_GetToken(t *testing.T) {
	cfg, _ := newConfigWithValues("testsecret", defaultTTL)
	gen := &Generator{cfg: cfg}

	tokenStr, err := gen.GetToken(123)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("testsecret"), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, float64(123), claims["user_id"])
}

func TestGenerator_ParseToken(t *testing.T) {
	cfg, _ := newConfigWithValues("testsecret", defaultTTL)
	gen := &Generator{cfg: cfg}

	tokenStr, _ := gen.GetToken(123)
	token, err := gen.ParseToken(tokenStr)
	require.NoError(t, err)
	assert.NotNil(t, token)
	assert.True(t, token.Valid)

	_, err = gen.ParseToken("invalid.token.here")
	assert.Error(t, err)

	wrongGen, _ := newConfigWithValues("wrongsecret", defaultTTL)
	wrongGen2 := &Generator{cfg: wrongGen}
	_, err = wrongGen2.ParseToken(tokenStr)
	assert.Error(t, err)
}
