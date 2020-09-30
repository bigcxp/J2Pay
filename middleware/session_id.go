package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"j2pay-server/pkg/setting"
)

const SessionIdName = "GO-SESSION-ID"

//手动实现一个session
func MakeSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		// request
		sessionId, _ := c.Cookie(SessionIdName)
		if sessionId == "" {
			sessionId = uuid.NewV4().String()
			c.SetCookie(SessionIdName, sessionId, 0, "/", setting.ApplicationConf.Domain, false, true)
		}
		c.Set("session_id", sessionId)
		c.Next()
	}
}
