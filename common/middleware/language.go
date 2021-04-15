package middleware

import (
	"strings"

	"github.com/dxasu/gostar/config"

	"github.com/gin-gonic/gin"
)

func SetLanguage(c *gin.Context) {
	var lang string
	if lang = c.Param("lang"); lang == "" {
		lang = c.Query("lang")
	}
	if len(lang) < 2 {
		lang = "en"
	}

	lang = strings.ToLower(lang)[0:2]
	langinfo := config.GetMapBYKey("language." + lang)

	c.Set("langinfo", langinfo)
}
