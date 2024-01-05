package ginweb

import (
	"fmt"
	"local/dbPool"
	"local/tools"

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

type student struct {
	name string
	age  int
}

func handleGetName(c *gin.Context) {
	var name = c.Param("name")
	var password = c.Param("password")
	fmt.Println("name: ", name)
	fmt.Println("password: ", password)
	return
}

func loginCheck(c *gin.Context) {
	var name = c.Param("name")
	var password = c.Param("password")
	var isPass = false

	defer func() {
		if isPass {
			fmt.Println("pass")
			c.Next()
		} else {
			fmt.Println("noPass")
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
			c.Next()
		}
	}
	return
}

func setLogInGroup(engine *gin.Engine) {
	loginGroup := engine.Group(loginPath)
	loginGroup.Use(loginCheck)
	loginGroup.GET("/:name/:password", handleGetName)
}

func InitGroup(engine *gin.Engine) {
	setLogInGroup(engine)
}
