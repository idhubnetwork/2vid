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

func CreateCredential(c *gin.Context, jt *jsontokens.JsonToken) {
	did, ok := jt.Get("did").(string)
	if !ok || len(did) != 32 {
		c.JSON(http.StatusForbidden, ActionErr{"jsontoken invalid or non did"})
	}

	var jwt JWT
	err := c.ShouldBind(&jwt)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{"invalid or non jwt to create"})
	}

	credentiual, err := db_mysql.VerifyWritedData(did, jwt.JsonWebToken)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
	}
}
