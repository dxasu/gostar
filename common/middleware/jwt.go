package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func JwtApiMiddleware(c *gin.Context) {
	token := c.GetHeader("token")
	if token == "" {
		token = c.Query("token")
	}
	userinfo := make(map[string]interface{})
	if userinfo == nil || userinfo["name"] == nil || userinfo["create_time"] == nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "验证失败",
		})
		c.Abort()
		return
	}
	createTime := int64(userinfo["create_time"].(float64))
	var expire int64 = 24 * 60 * 60
	nowTime := time.Now().Unix()
	if (nowTime - createTime) >= expire {
		c.JSON(200, gin.H{
			"code": 401,
			"msg":  "token失效",
		})
		c.Abort()
	}
	c.Set("user", userinfo["name"])
	//log.Println(userinfo)
	//if userinfo["type"]=="kefu"{
	c.Set("kefu_id", userinfo["kefu_id"])
	c.Set("kefu_name", userinfo["name"])
	c.Set("role_id", userinfo["role_id"])
	//}
}
