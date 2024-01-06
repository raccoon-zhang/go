package ginweb

import (
	"fmt"
	"local/dbPool"
	"local/tools"
	"net/http"

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

func handleGetName(c *gin.Context) {
	var name = c.PostForm("name")
	var password = c.PostForm("password")
	fmt.Println("name: ", name)
	fmt.Println("password: ", password)
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

	db, err := pool.NewDb()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		pool.DeleteDb(db)
	}()

	var sqlString = "select * from student where name = ?"
	ret, err := db.Query(sqlString, name)
	if err != nil {
		fmt.Println(err)
		return
	}
	for ret.Next() {
		var name string
		var age int
		var passwordHash string
		ret.Scan(&name, &age, &passwordHash)
		if tools.PasswordDecrypt(passwordHash, password) {
			isPass = true
		}
	}
}

func setLogInGroup(engine *gin.Engine) {
	loginGroup := engine.Group(loginPath)
	loginGroup.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", "")
	})
	loginGroup.POST("/check", loginCheck, handleGetName)
}

func InitGroup(engine *gin.Engine) {
	setLogInGroup(engine)
}
