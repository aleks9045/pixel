package api

import (
	"github.com/aleks9045/pixels/internal/service/user"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(e *gin.RouterGroup, userService *user.UserService, middleware gin.HandlerFunc) {

	userGroup := e.Group("/users")
	userGroup.Use(middleware)

	userGroup.GET("/me", userService.GetMe)
	userGroup.GET("/all", userService.GetAll)
	userGroup.GET("/:user_id", userService.GetUser)
	userGroup.PATCH("/:user_id", userService.PatchUser)
}
