package api

import (
	"github.com/aleks9045/pixels/internal/service/auth"
	"github.com/gin-gonic/gin"
)

func RegisterAuthGroup(e *gin.RouterGroup, service *auth.AuthService, middleware gin.HandlerFunc) {

	authGroup := e.Group("/auth")
	authGroup.Use(middleware)

	authGroup.POST("/logout", service.Logout)

}
