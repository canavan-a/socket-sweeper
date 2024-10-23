package main

import (
	"fmt"
	"main/sweeper"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Change this to match your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	r.SetTrustedProxies(nil)

	r.LoadHTMLGlob("dist/*.html")

	r.Static("/assets", "./assets")

	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"response": "server is running"})
	// })

	serverG := sweeper.InitServerGames()

	r.GET("/publish", serverG.PublisherRoute)
	r.GET("/subscribe", serverG.SubscriberRoute)
	r.GET("/menu", serverG.MenuConnectionRoute)

	r.NoRoute(func(c *gin.Context) {
		fmt.Println("you are in no route")
		c.HTML(200, "index.html", nil)
	})

	fmt.Println("server has started")
	r.Run(":80")
}
