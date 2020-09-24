package jwt

import (
	"log"
	"net/http"
	"set-flags/src/pkg/e"
	"set-flags/src/pkg/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// JWT jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.Request.Header.Get("Authorization")
		if token == "" || token == "null" || len(strings.Fields(token)) != 2 {
			code = e.INVALID_PARAMS
		} else {
			authToken := strings.Fields(token)[1]
			claims, err := utils.ParseToken(authToken)
			if err != nil {
				log.Printf("parse token failed, err: %v\n", err)
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
