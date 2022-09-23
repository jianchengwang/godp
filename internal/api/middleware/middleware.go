package middleware

import (
	"github.com/gin-gonic/gin"
	"godp/internal/config"
	"godp/internal/global"
	"godp/pkg/errorcode"
	"godp/pkg/permission"
	"godp/pkg/response"
	"net/http"
	"strings"
	"time"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", config.Config.App.CrossOrigin) // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		checkError := errorcode.Success
		token := c.GetHeader(global.JwtToken)
		if token == "" {
			token = c.Query(global.JwtToken)
		}
		token = strings.Replace(token, "Bearer ", "", 1)
		if token == "" {
			checkError = errorcode.ErrInvalidParam
		} else {
			claims, err := permission.ParseToken(token)
			if err != nil {
				checkError = errorcode.ErrInvalidToken
			} else if time.Now().Unix() > claims.ExpiresAt {
				checkError = errorcode.ErrTokenTimeout
			}
		}

		if checkError != errorcode.Success {
			response.Error(c, checkError)
			c.Abort()
			return
		}
		c.Next()
	}
}
