package middleware

import (
	"errors"
	"net/http"

	"github.com/aleks9045/pixels/internal/repository"
	"github.com/aleks9045/pixels/internal/service/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	accessCookieName  = "access_token"
	refreshCookieName = "refresh_token"
)

type AuthMiddleware struct {
	userRepository *repository.UserRepository
	tokenManager   *auth.TokenManager
}

func NewAuthMiddleware(userRepository *repository.UserRepository, tokenManager *auth.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		userRepository: userRepository,
		tokenManager:   tokenManager,
	}
}

func (a *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, accessErr := c.Cookie(accessCookieName)

		var userId string

		if accessErr == nil {
			userId, accessErr = a.tokenManager.ParseAccessToken(accessToken)
		}

		if accessErr != nil && !errors.Is(accessErr, jwt.ErrTokenExpired) && accessErr != http.ErrNoCookie {
			auth.ClearCookies(c)
			c.String(http.StatusUnauthorized, "unauthorized")
			return
		}

		if accessErr != nil {
			refreshToken, refreshErr := c.Cookie(refreshCookieName)
			if refreshErr != nil {
				auth.ClearCookies(c)
				c.String(http.StatusUnauthorized, "unauthorized")
				return
			}

			userId, refreshErr = a.tokenManager.ParseRefreshToken(refreshToken)
			if refreshErr != nil {
				auth.ClearCookies(c)
				c.String(http.StatusUnauthorized, "unauthorized")
				return
			}

			newAccessToken, refreshErr := a.tokenManager.CreateAccessToken(userId)
			if refreshErr != nil {
				c.String(http.StatusUnauthorized, "unauthorized")
				return
			}

			auth.SetAccessCookie(c, newAccessToken)
		}

		isAdmin, err := a.userRepository.GetIsAdmin(userId)
		if err != nil {
			auth.ClearCookies(c)
			c.String(http.StatusUnauthorized, "unauthorized")
			return
		}

		c.Set("user_id", userId)
		c.Set("is_admin", isAdmin)
		c.Next()
	}
}
