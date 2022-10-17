package middleware

import (
	"seckill/authentication/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth () gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的路径
		path := c.Request.URL.Path

		// 过滤掉不需要验证的接口
		if path == "/user/login" || path == "/user/register" || path == "/user/logout" {
			c.Next()
			return
		}
		// 从请求头中获取token
		authorization := c.Request.Header.Get("authorization")
		if authorization == "" {
			c.JSON(200, gin.H{
				"code": 401,
				"msg": "token为空",
			})
			// 验证失败，不再调用后续的函数处理
			c.Abort()
			return
		}
		parts := strings.Split(authorization, " ")
		if parts[0] != "Bearer" || len(parts) != 2 {
			c.JSON(200, gin.H{
				"code": 401,
				"msg": "token格式错误",
			})
			c.Abort()
			return
		}
		// 验证token
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(200, gin.H{
				"code": 401,
				"msg": "token无效",
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)

		c.Next()

	}
}