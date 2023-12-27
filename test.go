package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

import _ "github.com/go-sql-driver/mysql"

type student struct {
	name string
	age  int
}

func handleHome(c *gin.Context) {
	fmt.Println("hello,world")
}

func handleGetName(c *gin.Context) {
	var name = c.Param("name")
	fmt.Println(name)
}

func main() {
	db, err := sql.Open("mysql", "root:@/school")
	if err != nil {
		fmt.Println(err)
		return
	}

	ret, err := db.Query("select * from student")
	if err != nil {
		fmt.Println(err)
		return
	}
	print(ret)
	engine := gin.Default()
	engine.Use(func(c *gin.Context) {
		fmt.Println("middle1")
		c.Next()
	})
	engine.GET("/:name", handleGetName)
	// engine.Run(":8080")
}
