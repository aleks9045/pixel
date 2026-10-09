package auth

import (
	"github.com/gin-gonic/gin"
)

func GetUserIdFromCtx(ctx *gin.Context) (string, bool) {
	userId, ok := ctx.Get("user_id")
	if !ok {
		return "", ok
	}

	return userId.(string), ok
}

func GetIsAdminFromCtx(ctx *gin.Context) (bool, bool) {
	userId, ok := ctx.Get("is_admin")
	if !ok {
		return false, ok
	}

	return userId.(bool), ok
}
