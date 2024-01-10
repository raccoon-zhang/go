package main

import (
	"local/ginweb"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	engine := gin.Default()

	ginweb.InitEngine(engine)
	engine.Run(":8080")
}
