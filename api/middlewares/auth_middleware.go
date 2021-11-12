package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func AuthMiddleware(router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		if auth := viper.Get("server.http.basic_auth"); auth != nil {
			segment := strings.Split(auth.(string), ":")
			router.Use(gin.BasicAuth(gin.Accounts{
				segment[0]: segment[1],
			}))
		}
	}
}