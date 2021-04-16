# *goon*

## gin
* router
```
	router := gin.Default()
	router.LoadHTMLGlob("html/*.html")
	router.StaticFS("/static", http.Dir("./static"))

	router.GET("/form", middleware.SetLanguage, func(c *gin.Context) {
		common.CommitInfo(c, "", []common.PropInfo{
			{"Name:", "text", "name", "uname"},
			{"Email:", "text", "email", "email"},
			{"Document:", "text", "document", "123"},
		})
	})
```

* getValue
```
c.Param("uuid") //取得URL中参数  router.POST("/post/:uuid", func(c *gin.Context) {
c.Query("id") //查询请求URL后面的参数
c.PostForm("name") //从表单中查询参数
c.Request.Header.Get("User-Agent") //从header中读取
c.Request.PostFormValue("name") //return POST and PUT body parameters, URL query parameters are ignored
c.Get("current_manager") //从用户上下文读取值  
sessions.Default(c).get("user") //从session中读取值
```

## gorm 
* [gorm](https://gorm.io/zh_CN/docs/query.html)
    > https://gorm.io/zh_CN/docs/query.html

### gf 
* [goframe](https://github.com/gogf/gf)
    > https://github.com/gogf/gf

* ```	payData := &DB{}
	payData.InitDB("mysqlPay")
	payData.AutoMigrate(PersonMIlk{}, User{})
	count := 0
	payData.Create(PersonMIlk{21, "2", "3", "4"})
	payData.Create(PersonMIlk{30, "2", "3", "4"})```

```
stmt := db.Session(&Session{DryRun: true}).First(&user, 1).Statement
stmt.SQL.String() //=> SELECT * FROM `users` WHERE `id` = $1 ORDER BY `id`
stmt.Vars         //=> []interface{}{1}
```


```
func main() {
	defer log.Flush()

	config.Init()
	//go grpclib.GrpcInit()

	paylogic.Init()
	orm.Test()

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	log.Warningf("main signal:%+v", <-chSig)
}
```

---
*goon*
**goon**
~~goon~~


***
# other
1. [markdown网址](https://www.cnblogs.com/liugang-vip/p/6337580.html)
2. git tags
git tag -a "project-v0.6.1" -m "release go star"
git push --follow-tags
git push origin project -u --follow-tags
3. go build
go env -w GOPROXY=http://mirrors.aliyun.com/goproxy/