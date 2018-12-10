package handler

import (
	"2vid/mysql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idhubnetwork/jsontokens"
)

const (
	// 0011 0000
	DEFAULT_STATUS = 0x30
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

	credential, err := db_mysql.VerifyWritedData(did, jwt)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		return
	}

	err = db_mysql.CreateCredential(credential)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		return
	}

	c.JSON(http.StatusOK, ActionSuccess{"credential create successed"})
}
