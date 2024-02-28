package ginweb

import (
	"context"
	"fmt"
	"gptChat"
	"io"
	"local/dbPool"
	"local/tools"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/otiai10/gosseract/v2"
	"github.com/redis/go-redis/v9"
)

const loginPath string = "/login"
const registerPath string = "/register"
const chatPath string = "/chat"

var pool *dbPool.Pool
var redisCtx context.Context
var gptClients sync.Map

func init() {
	var err error
	var falioverOps = &redis.FailoverOptions{
		MasterName:    "localmaster",
		SentinelAddrs: []string{"localhost:26379", "localhost:26380"},
	}

	pool, err = dbPool.InitPool("mysql", "root:@/gptweb", falioverOps, 10)
	if err != nil {
		fmt.Println(err)
	}
	redisCtx = context.Background()
}

func savePage(c *gin.Context) {
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

func loginStatus(c *gin.Context) {
	//这里用来设置之后需要添加的右上角登陆状态
	id := sessions.Default(c).Get("userKey")
	if id == nil {
		c.JSON(http.StatusOK, gin.H{"isLogin": "false"})
	} else {
		c.JSON(http.StatusOK, gin.H{"isLogin": "true"})
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
	rdbRead := pool.NewRedisCliForRead()

	defer pool.DeleteRedisCli(rdbRead)

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
	defer pool.DeleteDb(db)
	var sqlString = "select name,password from user where name = ?"
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
			defer pool.DeleteRedisCli(rdbWrite)
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

// 给前端用的，至于延迟多少时间前端自己设置，后端不处理，这里指返回之前保存的页面
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

func setRegisterInGroup(engine *gin.Engine) {
	registerGroup := engine.Group(registerPath)
	registerGroup.GET("/", multiLoginCheck, func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", "")
	})
	registerGroup.POST("/", registerUser)
}

func setChatGroup(engine *gin.Engine) {
	chatGroup := engine.Group(chatPath, mustLogin)
	chatGroup.GET("/", savePage, func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", "")
	})
	chatGroup.POST("/queryGpt", func(c *gin.Context) {
		msg := c.PostForm("userMessage")
		responseChan := make(chan interface{})
		go func() {
			var client interface{}
			if value, ok := gptClients.Load(sessions.Default(c).Get("userKey")); !ok {
				client = gptChat.DefaultClient()
				gptClients.Store(sessions.Default(c).Get("userKey"), client)
			} else {
				client = value
			}
			if cli, ok := client.(gptChat.LocalClient); ok {
				data, err := cli.QueryGpt(msg)
				if err != nil {
					fmt.Println(err)
					responseChan <- "something wrong, not your fault"
				} else {
					responseChan <- data
					fmt.Println(data)
				}
			}
			defer close(responseChan)
		}()

		select {
		case data := <-responseChan:
			c.JSON(http.StatusOK, gin.H{"botResponce": data})
		case <-time.After(time.Second * 30): // 设置超时时间为30秒
			c.JSON(http.StatusOK, gin.H{"botResponce": "Gpt Operation Timed Out"})
		}
	})

	chatGroup.POST("/ocr", func(c *gin.Context) {
		image, err := c.FormFile("image")
		if err != nil {
			fmt.Println(err)
			return
		}
		imageReader, err := image.Open()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer imageReader.Close()
		imageData, err := io.ReadAll(imageReader)
		if err != nil {
			fmt.Println(err)
			return
		}
		langs, err := gosseract.GetAvailableLanguages()
		if err != nil {
			fmt.Println(err)
			return
		}
		ocrCli := gosseract.NewClient()
		defer ocrCli.Close()
		//将现有语言支持都加入
		ocrCli.SetLanguage(langs...)
		ocrCli.SetImageFromBytes(imageData)
		text, _ := ocrCli.Text()
		c.JSON(http.StatusOK, gin.H{"text": text})
	})
}

func registerUser(c *gin.Context) {
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

	if sqlRegisterUser(name, age, password) {
		redisRegisterUser(name, password)
		isPass = true
	}
}

func redisRegisterUser(name, password string) bool {
	passwordHash, err := tools.PasswordEncrypt(password)
	if err != nil {
		fmt.Println(err)
		return false
	}
	rdbWrite := pool.NewRedisCliForWrite()
	defer pool.DeleteRedisCli(rdbWrite)

	err = rdbWrite.Set(redisCtx, name, passwordHash, time.Hour*24).Err()
	if err != nil {
		fmt.Println(err)
		return false
	} else {
		return true
	}
}

func sqlRegisterUser(name, age, password string) bool {
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

	defer pool.DeleteDb(db)

	stmt, err := db.Prepare("insert into user(name,age,password) values(?,?,?)")

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
	engine.Static("/static", "../templates")
	engine.LoadHTMLGlob("../templates/*.html")
}

func initRouter(engine *gin.Engine) {
	//--------------------页面渲染---------------------------
	//首页
	engine.GET("/", savePage, func(c *gin.Context) {
		fmt.Println("hello world")
		//TODO: 设计一个首页,暂时使用聊天页面代替
		c.Redirect(http.StatusTemporaryRedirect, chatPath)
	})

	//登陆注册
	setLogInGroup(engine)
	setRegisterInGroup(engine)
	engine.GET("/logout", func(c *gin.Context) {
		removeSessionVal("userKey", c)
		c.Redirect(http.StatusTemporaryRedirect, loginPath)
	})

	//聊天页
	setChatGroup(engine)

	//404页面
	engine.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{"domain": "127.0.0.1:8080"})
	})
	//------------------------------------------------------

	//--------------------问询接口---------------------------
	//前页面查询
	engine.GET("/backPage", backPage)

	//登陆状态检查
	engine.GET("/loginStatus", loginStatus)
	//------------------------------------------------------
}
