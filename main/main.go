package main

import (
	"local/ginweb"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	ginweb.InitEngine(engine)
	engine.Run(":8080")
}
