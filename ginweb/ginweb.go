package ginweb

import (
	"context"
	"fmt"
	"local/dbPool"
	"local/tools"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const loginPath string = "/login"
const registePath string = "/registe"

var pool *dbPool.Pool
var redisCtx context.Context
var rdbRead *redis.Client
var rdbWrite *redis.Client

func init() {
	var err error
	pool, err = dbPool.InitPool("mysql", "root:@/school", 10)
	if err != nil {
		fmt.Println(err)
	}
	redisCtx = context.Background()
	rdbRead = redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "",
		DB:       0,
	})
	rdbWrite = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func savePrePage(c *gin.Context) {
	//其他的中间件执行完毕之后再更新当前页面，否则不更新
	c.Next()
	setSessionVal("prepage", c.Request.URL.Path, c)
	fmt.Println("save prepage:", c.Request.URL.Path)
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

func setSessionVal(key, value any, c *gin.Context) {
	session := sessions.Default(c)
	session.Set(key, value)
	session.Save()
}

func removeSessionVal(key any, c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(key)
	session.Save()
}

func redisloginCheck(name, password string, c *gin.Context) bool {
	passwordHash, err := rdbRead.Get(redisCtx, name).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	if tools.PasswordDecrypt(passwordHash, password) {
		setSessionVal("userKey", passwordHash, c)
		fmt.Println("redis check")
		return true
	}
	return false
}

func sqlLoginCheck(name, password string, c *gin.Context) bool {
	db, err := pool.NewDb()
	if err != nil {
		return false
	}
	defer func() {
		pool.DeleteDb(db)
	}()
	var sqlString = "select name,password from student where name = ?"
	ret, err := db.Query(sqlString, name)
	if err != nil {
		fmt.Println(err)
		return false
	}
	for ret.Next() {
		var name string
		var passwordHash string
		ret.Scan(&name, &passwordHash)
		if tools.PasswordDecrypt(passwordHash, password) {
			setSessionVal("userKey", name, c)
			err := rdbWrite.Set(redisCtx, name, passwordHash, time.Hour*24).Err()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("sql check")
			return true
		}
	}
	return false
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
	}()
	if redisloginCheck(name, password, c) {
		fmt.Println("redis check")
		isPass = true
	} else if sqlLoginCheck(name, password, c) {
		fmt.Println("sql check")
		isPass = true
	} else {
		c.Abort()
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
	loginGroup.POST("/", loginCheck)
}

func setRegisteInGroup(engine *gin.Engine) {
	registeGroup := engine.Group(registePath)
	registeGroup.GET("/", multiLoginCheck, func(c *gin.Context) {
		c.HTML(http.StatusOK, "registe.html", "")
	})
	registeGroup.POST("/", registeUser)
}

func registeUser(c *gin.Context) {
	var name string
	var password string
	var age string
	if c.Request.Method == http.MethodGet {
		name = c.Param("name")
		password = c.Param("password")
		age = c.Param("age")
	} else if c.Request.Method == http.MethodPost {
		name = c.PostForm("name")
		password = c.PostForm("password")
		age = c.PostForm("age")
	}

	var isPass = false
	defer func() {
		if isPass {
			fmt.Println("registe success")
			c.JSON(http.StatusOK, gin.H{"status": "true"})
			c.Next()
		} else {
			fmt.Println("registe failed")
			c.JSON(http.StatusUnauthorized, gin.H{"status": "false"})
			c.Abort()
		}
	}()

	if sqlRegisteUser(name, age, password) {
		redisRegisteUser(name, password)
		isPass = true
	}
}

func redisRegisteUser(name, password string) bool {
	passwordHash, err := tools.PasswordEncrypt(password)
	if err != nil {
		fmt.Println(err)
		return false
	}
	err = rdbWrite.Set(redisCtx, name, passwordHash, time.Hour*24).Err()
	if err != nil {
		fmt.Println(err)
		return false
	} else {
		return true
	}
}

func sqlRegisteUser(name, age, password string) bool {
	passwordHash, err := tools.PasswordEncrypt(password)

	if err != nil {
		fmt.Println(err)
		return false
	}

	db, err := pool.NewDb()

	if err != nil {
		fmt.Println(err)
		return false
	}

	defer func() {
		pool.DeleteDb(db)
	}()

	stmt, err := db.Prepare("insert into student(name,age,password) values(?,?,?)")

	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = stmt.Exec(name, age, passwordHash)

	if err != nil {
		fmt.Println(err)
		return false
	} else {
		return true
	}
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
	setRegisteInGroup(engine)
	engine.GET("/backPage", backPage)
	//保存界面
	engine.Use(savePrePage)
	engine.GET("/", func(c *gin.Context) {
		fmt.Println("hello world")
	})
	engine.GET("/check", mustLogin, func(c *gin.Context) {

	})
	engine.GET("/logout", func(c *gin.Context) {
		removeSessionVal("userKey", c)
	})
}
