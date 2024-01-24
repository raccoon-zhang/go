package ginweb

import (
	"context"
	"fmt"
	"gptChat"
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
const chatPath string = "/chat"

var pool *dbPool.Pool
var redisCtx context.Context

func init() {
	var err error
	var falioverOps = &redis.FailoverOptions{
		MasterName:    "localmaster",
		SentinelAddrs: []string{"localhost:26379", "localhost:26380"},
	}

	pool, err = dbPool.InitPool("mysql", "root:@/school", falioverOps, 10)
	if err != nil {
		fmt.Println(err)
	}
	redisCtx = context.Background()
}

func savePrePage(c *gin.Context) {
	//其他的中间件执行完毕之后再更新当前页面，否则不更新
	c.Next()
	setSessionVal("prepage", c.Request.URL.Path, c)
}

func mustLogin(c *gin.Context) {
	id := sessions.Default(c).Get("userKey")
	if id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}

func sessionCheck(c *gin.Context) {
	//这里用来设置之后需要添加的右上角登陆状态
	id := sessions.Default(c).Get("userKey")
	if id == nil {
		//没有登陆
	} else {
		//已经登陆
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
	rdbRead := pool.NewRedisCliForRead(&redis.Options{
		Addr:     "localhost:6380",
		Password: "",
		DB:       0,
	})

	defer func() {
		pool.DeleteRedisCli(rdbRead)
	}()

	passwordHash, err := rdbRead.Get(redisCtx, name).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	if tools.PasswordDecrypt(passwordHash, password) {
		setSessionVal("userKey", name, c)
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
			rdbWrite := pool.NewRedisCliForWrite()
			defer func() {
				pool.DeleteRedisCli(rdbWrite)
			}()
			err := rdbWrite.Set(redisCtx, name, passwordHash, time.Hour*24).Err()
			if err != nil {
				fmt.Println(err)
			}
			return true
		}
	}
	return false
}

func loginCheck(c *gin.Context) {
	var name string
	var password string
	name = c.PostForm("name")
	password = c.PostForm("password")

	var isPass = false
	defer func() {
		if isPass {
			c.JSON(http.StatusOK, gin.H{"status": "true"})
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "false"})
			c.Abort()
		}
	}()
	if redisloginCheck(name, password, c) {
		isPass = true
	} else if sqlLoginCheck(name, password, c) {
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

func setChatGrop(engine *gin.Engine) {
	chatGroup := engine.Group(chatPath, mustLogin)
	chatGroup.GET("/", savePrePage, func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", "")
	})
	chatGroup.POST("/queryGpt", func(c *gin.Context) {
		msg := c.PostForm("userMessage")
		responceChan := make(chan interface{})
		go func() {
			data, err := gptChat.QueryGpt(msg)
			if err != nil {
				fmt.Println(err)
				responceChan <- "something wrong, not your fault"
			} else {
				responceChan <- data
				fmt.Println(data)
			}
			defer close(responceChan)
		}()

		select {
		case data, ok := <-responceChan:
			if ok {
				c.JSON(http.StatusOK, gin.H{"botResponce": data})
			}
		case <-time.After(time.Second * 2): // 设置超时时间为2秒
			c.JSON(http.StatusOK, gin.H{"botResponce": "Gpt Operation Timed Out"})
		}
	})
}

func registeUser(c *gin.Context) {
	var name string
	var password string
	var age string
	name = c.PostForm("name")
	password = c.PostForm("password")
	age = c.PostForm("age")

	var isPass = false
	defer func() {
		if isPass {
			c.JSON(http.StatusOK, gin.H{"status": "true"})
			c.Next()
		} else {
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
	rdbWrite := pool.NewRedisCliForWrite()
	defer func() {
		pool.DeleteRedisCli(rdbWrite)
	}()

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

	//首页
	engine.GET("/", savePrePage, func(c *gin.Context) {
		fmt.Println("hello world")
	})

	//登陆注册
	setLogInGroup(engine)
	setRegisteInGroup(engine)
	engine.GET("/logout", func(c *gin.Context) {
		removeSessionVal("userKey", c)
	})
	engine.GET("/backPage", backPage)

	//聊天页
	setChatGrop(engine)
}
