package main

import "github.com/gin-gonic/gin"

func main() {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"response": "server is running"})
	})

	r.Run(":80")
}
