package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"

	"local/dbPool"
)

var pool *dbPool.Pool

func init() {
	var err error
	pool,err = dbPool.InitPool("mysql", "root:@/school", 10)
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
	fmt.Println(name)
}

func loginCheck(c *gin.Context) {
	var name = c.Param("name")
	db,err := pool.NewDb()
	if err != nil {
		fmt.Println(err)
		return
	}
	var sqlString = "select * from student where name = ?"
	ret, err := db.Query(sqlString, name)
	if err != nil {
		fmt.Println(err)
		return
	}
	for ret.Next() {
		var name string
		var age int
		ret.Scan(&name,&age)
		if age == 24 {
			fmt.Println("pass")
			c.Next()
		} else {
			fmt.Println("nopass ","age:",age)
			return
		}
	}
}

func main() {
	engine := gin.Default()
	engine.Use(loginCheck)
	engine.GET("/:name", handleGetName)
	engine.Run(":8080")
	pool.DestroyPool()
}
