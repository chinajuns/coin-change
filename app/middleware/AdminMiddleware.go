package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okc/utils"
)

// AdminMiddleware
// admin路由中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		lang := c.GetHeader("lang")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"type":    "999",
				"message": utils.GetLangMessage(lang, utils.UserAuthError),
			})
			c.Abort()
			return
		}
		userClaims, err := utils.AnalyseToken(token)
		//log.Println("userClaims :", userClaims)
		if err != nil {
			_ = utils.WriteErrorLog(fmt.Sprintf("ApiMiddleware [ERROR] %s \n", err))
			c.JSON(http.StatusOK, gin.H{
				"type":    "999",
				"message": utils.GetLangMessage(lang, utils.UserAuthError),
			})
			c.Abort()
			return
		}
		c.Set("lang", lang)
		c.Set("user_info", userClaims)
		c.Next()
	}
}
