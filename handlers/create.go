package handler

import (
	"2vid/logger"
	"2vid/mysql"
	"2vid/redis"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idhubnetwork/jsontokens"
)

// Create a new Credential, param binding a JWT
func createCredential(c *gin.Context, jt *jsontokens.JsonToken) {
	did, ok := jt.Get("did").(string)
	if !ok || len(did) != 52 {
		c.JSON(http.StatusForbidden, ActionErr{"jsontoken invalid or non did"})
		return
	}

	tmp, ok := c.Get("credential")
	jwt, ok := tmp.(string)
	if !ok {
		c.JSON(http.StatusForbidden, ActionErr{"invalid or non jwt to create"})
		return
	}

	_, err := db_mysql.VerifyWritedData(did, jwt)
	if err != nil {
		logger.Log.Warn(err)
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		return
	}

	err = db_redis.Publish("create", 0, 0, jwt)
	if err != nil {
		logger.Log.Error(err)
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		return
	}

	c.JSON(http.StatusOK, ActionSuccess{"credential create successed"})
}
