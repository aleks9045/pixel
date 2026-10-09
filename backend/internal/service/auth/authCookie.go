package auth

import "github.com/gin-gonic/gin"

func SetAccessCookie(c *gin.Context, token string) {
	c.SetCookie(
		accessCookieName,
		token,
		int(accessCookieTTL.Seconds()),
		"/",
		"",
		false,
		true,
	)
}

func SetRefreshCookie(c *gin.Context, token string) {
	c.SetCookie(
		refreshCookieName,
		token,
		int(refreshCookieTTL.Seconds()),
		"/",
		"",
		false,
		true,
	)
}

func ClearCookies(c *gin.Context) {
	c.SetCookie(
		accessCookieName,
		"",
		0,
		"/",
		"",
		true,
		true,
	)

	c.SetCookie(
		refreshCookieName,
		"",
		0,
		"/auth",
		"",
		false,
		true,
	)
}
