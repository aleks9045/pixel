package api

import (
	"github.com/aleks9045/pixels/internal/service/pixel"
	"github.com/gin-gonic/gin"
)

func RegisterPixelRoutes(r *gin.RouterGroup, service *pixel.PixelService, middleware gin.HandlerFunc) {

	pixelGroup := r.Group("/pixel")
	pixelGroup.Use(middleware)

	pixelGroup.GET("/field", service.GetField)
	pixelGroup.GET("/author", service.GetAuthor)

}
