package auth

import (
	"time"

	"github.com/aleks9045/pixels/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour

	accessTokenType  = "access"
	refreshTokenType = "refresh"
)

type Claims struct {
	Type string `json:"type"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secret []byte
}

func NewTokenManager(cfg *config.Config) *TokenManager {

	return &TokenManager{
		secret: []byte(cfg.JwtSecret),
	}
}

func (t *TokenManager) createToken(
	userID string,
	tokenType string,
	ttl time.Duration,
) (string, error) {
	claims := Claims{
		Type:      tokenType,
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(t.secret)
}

func (t *TokenManager) CreateAccessToken(userID string) (string, error) {
	return t.createToken(
		userID,
		accessTokenType,
		accessTokenTTL,
	)
}

func (t *TokenManager) CreateRefreshToken(userID string) (string, error) {
	return t.createToken(
		userID,
		refreshTokenType,
		refreshTokenTTL,
	)
}

func (t *TokenManager) ParseAccessToken(tokenString string) (string, error) {
	return t.parseToken(tokenString, accessTokenType)
}

func (t *TokenManager) ParseRefreshToken(tokenString string) (string, error) {
	return t.parseToken(tokenString, refreshTokenType)
}

func (t *TokenManager) parseToken(
	tokenString string,
	expectedType string) (string, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return t.secret, nil
		},
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Alg(),
		}),
	)

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", jwt.ErrTokenSignatureInvalid
	}

	if claims.Type != expectedType {
		return "", jwt.ErrInvalidType
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return "", jwt.ErrTokenExpired
	}

	return claims.Subject, nil
}
