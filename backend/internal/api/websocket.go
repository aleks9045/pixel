package api

import (
	"github.com/aleks9045/pixels/internal/service/ws"
	"github.com/gin-gonic/gin"
)

func RegisterWebsocketRoutes(r *gin.RouterGroup, service *ws.WsService, middleware gin.HandlerFunc) {

	wsGroup := r.Group("/ws")
	wsGroup.Use(middleware)

	wsGroup.GET("/connect", service.NewConnection)

}
