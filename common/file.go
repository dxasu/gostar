package common

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"

	log "github.com/dxasu/gostar/util/glog"

	"github.com/gin-gonic/gin"
)

func upfile(router *gin.Engine) {
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Static("/file", "./template")
	upload(router)
	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})
}

func cookie(router *gin.Engine) {
	router.GET("/cookie", func(c *gin.Context) {

		cookie, err := c.Cookie("gin_cookie")

		if err != nil {
			cookie = "NotSet"
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}

		log.Infof("Cookie value: %s \n", cookie)
	})
}

// IsFileExist 判断文件文件夹是否存在
func IsFileExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	}
	//我这里判断了如果是0也算不存在
	if fileInfo.Size() == 0 {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}

func upload(router *gin.Engine) {
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// single file
		file, _ := c.FormFile("file")
		log.Infoln(file.Filename)

		// Upload the file to specific dst.
		c.SaveUploadedFile(file, "./file")

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})
}

func multiUpload(router *gin.Engine) {
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/uploads", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Infoln(file.Filename)

			// Upload the file to specific dst.
			c.SaveUploadedFile(file, "./file")
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})
}

func downloadHandler(c *gin.Context) {
	fileName := "ceshi.png"  //filename  文件名
	path := "./data/images/" //文件存放目录防止用户下载其他目录文件
	file, err := os.Open(path + fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	fileNames := url.QueryEscape(fileName) // 防止中文乱码

	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+fileNames+"\"")

	if err != nil {
		fmt.Println("read file err:", err.Error())
		c.Status(http.StatusBadRequest)
	} else {
		c.String(http.StatusOK, string(content))
	}
}

func TranslateEn2Ch(text string) (string, error) {
	url := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=zh-cn&tl=en&dt=t&q=%s", url.QueryEscape(text))
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//返回的json反序列化比较麻烦, 直接字符串拆解
	ss := string(bs)
	ss = strings.ReplaceAll(ss, "[", "")
	ss = strings.ReplaceAll(ss, "]", "")
	ss = strings.ReplaceAll(ss, "null,", "")
	ss = strings.Trim(ss, `"`)
	ps := strings.Split(ss, `","`)
	return ps[0], nil
}

type student struct {
	Name string
	Age  int8
}

func templ(r *gin.Engine) {
	r.LoadHTMLGlob("tmpl/*")

	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/arr", func(c *gin.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gin.H{
			"title":  "Gin",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	t := template.Must(
		template.New("report").Parse("I am {{.Age}} years old.\n"))
	t.Execute(os.Stdout, struct{ Age int }{5})
	t.Execute(os.Stdout, map[string]int{"Age": 5})
}
