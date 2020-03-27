package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Session(secret string) gin.HandlerFunc {
	store := cookie.NewStore([]byte(secret))

	//
	store.Options(sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   7 * 86400,
		Secure:   false,
		HttpOnly: true,
		SameSite: 0,
	})
	return sessions.Sessions("gin-session", store)
}
