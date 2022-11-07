package middle

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/utils"
	"net/http"
)

//jwt中间件，用于验证token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//对登录接口放行
		if (len(c.Request.URL.String()) >= 10 && c.Request.URL.String()[0:10] == "/api/login"){
			c.Next()
		} else {
			token := c.Request.Header.Get("Authorization")
			if token == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "请求未携带token，无权限访问",
					"data": nil,
				})
				c.Abort()
				return
			}
			claims, err := utils.JWTToken.ParseToken(token)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": err.Error(),
					"data": nil,
				})
				c.Abort()
				return
			}
			c.Set("claims", claims)
		}
	}
}