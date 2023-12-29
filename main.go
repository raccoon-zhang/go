package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

var dbPool = make([]*sql.DB, 0)

func init() {
	fmt.Println("pool init")
	db, err := sql.Open("mysql", "root:@/school")
	if err != nil {
		fmt.Println(err)
		return
	}
	dbPool = append(dbPool, db)
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
	var db = dbPool[len(dbPool)-1]
	if db == nil {
		return
	}
	var sqlString = "select * from student where name = ?"
	var ret, err = db.Query(sqlString, name)
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
}
