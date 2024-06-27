package main

import (
	"fmt"
	"main/sweeper"

	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"response": "server is running"})
	})

	serverG := sweeper.InitServerGames()

	r.GET("/publish", serverG.PublisherRoute)
	r.GET("/subscribe", serverG.SubscriberRoute)
	r.GET("/menu", serverG.MenuConnectionRoute)

	fmt.Println("server has started")
	r.Run(":80")
}
