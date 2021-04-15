package middleware

import (
	log "github.com/dxasu/gostar/util/glog"

	"github.com/gin-gonic/gin"
)

var black []string

func init() {
	black = []string{
		"127.0.0.1",
		"111.1.2.3",
	}
}

func checkip(ip string) bool {
	for _, v := range black {
		if v == ip {
			return true
		}
	}
	return false
}

func Ipblack(c *gin.Context) {
	ip := c.ClientIP()
	if checkip(ip) {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "IP已被加入黑名单",
		})
		c.Abort()
		return
	}
	log.Infoln("client ip:", ip)
}
