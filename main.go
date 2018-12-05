package main

import (
	"2vid/handlers"
	"2vid/middleware"
	"2vid/mysql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	defer db_mysql.DB_mysql.Close()

	router := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// Simple group: v1
	v1 := router.Group("/api/v1", middleware.Authentication())
	{
		/*
			// URL EXAMPLE /credentials?iss=did:idhub:0x1234567890&aud=did:idhub:0987654321&sub=test
			v1.GET("/credentials", getCredential)
			v1.POST("/credentials", createCredential)
			v1.PUT("/credentials", updateCredential)
			v1.DELETE("/credentials", deleteCredential)

			// URL EXAMPLE /exceptions?iss=did:idhub:0x1234567890&aud=did:idhub:0987654321&sub=test
			v1.GET("/exceptions", recoverCredential)
			v1.POST("/exceptions", recoverCredential)
		*/
		v1.GET("/", handler.HandleCredential)
		v1.POST("/", handler.HandleCredential)
	}

	router.Run(":8000")
}
