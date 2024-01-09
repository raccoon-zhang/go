package main

import (
	"local/ginweb"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	engine := gin.Default()
	store := memstore.NewStore([]byte("user_id"))
	engine.Use(sessions.Sessions("user_id", store))

	engine.LoadHTMLGlob("../templates/*")
	ginweb.InitRouter(engine)
	engine.Run(":8080")
}
