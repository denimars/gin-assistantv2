package code

func Run() string {
	return `
package app

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run() {
	//db := db.Connection()
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(cors.Default())
	//api := router.Group("/api")
	//import your service in here...
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	router.Run(fmt.Sprintf(":%v", port))
}
	`
}
