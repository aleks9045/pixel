package api

import (
	"github.com/aleks9045/pixels/internal/service/oauth"
	"github.com/gin-gonic/gin"
)

func RegisterOauthRoutes(r *gin.RouterGroup, service *oauth.OauthService) {

	authGroup := r.Group("/oauth")

	authGroup.GET("/:oauth/login", service.Login)
	authGroup.GET("/:oauth/callback", service.Callback)

}
