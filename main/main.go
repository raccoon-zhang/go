package main

import (
	"github.com/gin-gonic/gin"
	"local/ginweb"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	engine := gin.Default()
	engine.LoadHTMLGlob("../templates/*")
	ginweb.InitGroup(engine)
	engine.Run(":8080")
}
