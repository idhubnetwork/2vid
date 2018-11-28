package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Simple group: v1
	v1 := router.Group("/api/v1")
	{
		// URL EXAMPLE /credentials?iss=did:idhub:0x1234567890&aud=did:idhub:0987654321&sub=test
		v1.GET("/credentials", getCredentials)
		v1.POST("/credentials", createCredentials)
		v1.PUT("/credentials", updateCredentials)
		v1.DELETE("/credentials", deleteCredentials)
	}

	router.Run(":8080")
}
