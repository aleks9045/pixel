package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	accessCookieName  = "access_token"
	refreshCookieName = "refresh_token"

	accessCookieTTL  = 30 * time.Minute
	refreshCookieTTL = 60 * 24 * time.Hour
)

type AuthService struct {
	tokenManager *TokenManager
}

func NewAuthService(tokenManager *TokenManager) *AuthService {
	return &AuthService{
		tokenManager: tokenManager,
	}
}

func (s *AuthService) Authenticate(c *gin.Context, userId *string) error {
	accessToken, err := s.tokenManager.CreateAccessToken(*userId)
	if err != nil {
		return err
	}

	refreshToken, err := s.tokenManager.CreateRefreshToken(*userId)
	if err != nil {
		return err
	}

	SetAccessCookie(c, accessToken)
	SetRefreshCookie(c, refreshToken)

	return nil
}

func (s *AuthService) Logout(c *gin.Context) {
	ClearCookies(c)

	c.String(http.StatusOK, "logged out")
}
