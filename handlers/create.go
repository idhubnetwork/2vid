package handler

import (
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

	tmp := jsontokens.NewJWT()
	err = tmp.SetJWT(jwt.JsonWebToken)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{"invalid jwt to create"})
	}
	err = tmp.Verify()
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{"invalid jwt signature"})
	}

	if did != tmp.Get("iss").(string) {
		c.JSON(http.StatusForbidden, ActionErr{"non jwt issuer can not create"})
	}
}
