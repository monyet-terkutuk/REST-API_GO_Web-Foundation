package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	service := NewService()
	userID := 123

	token, err := service.GenerateToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	service := NewService()
	userID := 123

	// Generate token
	token, err := service.GenerateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate token
	validToken, err := service.ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, validToken)
	assert.True(t, validToken.Valid)
	claims, ok := validToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, float64(userID), claims["user_id"])
}
