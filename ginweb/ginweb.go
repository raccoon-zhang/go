package ginweb

import (
	"fmt"
	"local/dbPool"
	"local/tools"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

const loginPath string = "/login"

var pool *dbPool.Pool

func init() {
	var err error
	pool, err = dbPool.InitPool("mysql", "root:@/school", 10)
	if err != nil {
		fmt.Println(err)
	}
}

func savePrePage(c *gin.Context) {
	//其他的中间件执行完毕之后再更新当前页面，否则不更新
	c.Next()
	session := sessions.Default(c)
	session.Set("prepage", c.Request.URL.Path)
	session.Save()
	fmt.Println("save prepage:", c.Request.URL.Path)
}
func handleGetName(c *gin.Context) {
	var name = c.PostForm("name")
	var password = c.PostForm("password")
	fmt.Println("name: ", name)
	fmt.Println("password: ", password)
}

func mustLogin(c *gin.Context) {
	id := sessions.Default(c).Get("userKey")
	if id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}

func sessionCheck(c *gin.Context) {
	id := sessions.Default(c).Get("userKey")
	if id == nil {
		fmt.Println("you have not login")
	} else {
		fmt.Println(id)
		fmt.Println("yong have login")
	}
}

func loginCheck(c *gin.Context) {
	var name string
	var password string
	if c.Request.Method == http.MethodGet {
		name = c.Param("name")
		password = c.Param("password")
	} else if c.Request.Method == http.MethodPost {
		name = c.PostForm("name")
		password = c.PostForm("password")
	}

	var isPass = false
	db, err := pool.NewDb()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if isPass {
			fmt.Println("pass")
			c.JSON(http.StatusOK, gin.H{"status": "true"})
			c.Next()
		} else {
			fmt.Println("noPass")
			c.JSON(http.StatusUnauthorized, gin.H{"status": "false"})
			c.Abort()
		}
		pool.DeleteDb(db)
	}()

	var sqlString = "select name,password from student where name = ?"
	ret, err := db.Query(sqlString, name)
	if err != nil {
		fmt.Println(err)
		return
	}
	for ret.Next() {
		var name string
		var passwordHash string
		ret.Scan(&name, &passwordHash)
		if tools.PasswordDecrypt(passwordHash, password) {
			isPass = true
			session := sessions.Default(c)
			session.Set("userKey", name)
			session.Save()
		}
	}
}

func getPrePageUrl(c *gin.Context) string {
	var url, ok = sessions.Default(c).Get("prepage").(string)
	if !ok {
		url = "/"
	}
	return url
}

// 给前端用的，至于延迟多少时间前端自己设置，后端不处理
func backPage(c *gin.Context) {
	var url = getPrePageUrl(c)
	c.JSON(http.StatusOK, gin.H{"prepage": url})
}

func multiLoginCheck(c *gin.Context) {
	id := sessions.Default(c).Get("userKey")
	if id != nil {
		var url = getPrePageUrl(c)
		fmt.Println("no need to login multiptly ,redirect to", url)
		c.Redirect(http.StatusTemporaryRedirect, url)
		c.Abort()
	}
}

func setLogInGroup(engine *gin.Engine) {
	loginGroup := engine.Group(loginPath)
	loginGroup.GET("/", multiLoginCheck, func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", "")
	})
	loginGroup.POST("/", loginCheck, handleGetName)
}

func InitEngine(engine *gin.Engine) {
	initSession(engine)
	initSource(engine)
	initRouter(engine)
}

func initSession(engine *gin.Engine) {
	store := memstore.NewStore([]byte("user_id"))
	engine.Use(sessions.Sessions("user_id", store))
}

func initSource(engine *gin.Engine) {
	engine.LoadHTMLGlob("../templates/*")
}

func initRouter(engine *gin.Engine) {
	//保存之前访问的页面，用于重复登陆返回页面,这里要先检查是否登陆在保存页面
	engine.Use(sessionCheck)
	//login界面和跳转界面不用保存，总不能登陆之后再跳回login或者跳转界面，会死循环
	setLogInGroup(engine)
	engine.GET("/backPage", backPage)
	//保存界面
	engine.Use(savePrePage)
	engine.GET("/", func(c *gin.Context) {
		fmt.Println("hello world")
	})
	engine.GET("/check", mustLogin, func(c *gin.Context) {

	})
}
