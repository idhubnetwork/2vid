package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/idhubnetwork/jsontokens"
)

const DESTINATION = "url"

type AuthErr struct {
	Authentication string `json:"FaliedAuthentication"`
}

type Token struct {
	Token string `form:"token" json:"token" xml:"token" binding:"required"`
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		tmp, err := searchToken(c)
		if err != nil {
			c.JSON(http.StatusForbidden, AuthErr{err.Error()})
		}
		jt := jsontokens.NewJsonToken()
		jt.SetToken(tmp)
		if err = jt.Verify(); err != nil {
			c.JSON(http.StatusForbidden, AuthErr{err.Error()})
		}
		expiration, ok := jt.Get("expiration").(int)
		if !ok {
			c.JSON(http.StatusForbidden, AuthErr{"jsontoken non expiration"})
		}
		if int64(expiration) > time.Now().Unix() {
			c.JSON(http.StatusForbidden, AuthErr{"jsontoken expired"})
		}
		destination, ok := jt.Get("destination").(string)
		if !ok {
			c.JSON(http.StatusForbidden, AuthErr{"jsontoken non destination"})
		}
		if destination != DESTINATION {
			c.JSON(http.StatusForbidden, AuthErr{"invalid access url destination"})
		}
		c.Set("DIDJsonToken", jt)
	}
}

func searchToken(c *gin.Context) (string, error) {
	tmp := c.Request.Header.Get("Authorization")
	var token Token
	if len(tmp) < 14 || tmp[0:13] != "DIDJsonToken " {
		if err := c.ShouldBind(&token); err != nil {
			return "", errors.New("non DID Json Token")
		}
		tmp = token.Token
		if len(tmp) < 14 || tmp[0:13] != "DIDJsonToken " {
			return "", errors.New("invalid DID Json Token")
		}
		return "", errors.New("invalid DID Json Token")
	}
	return string(tmp[13:]), nil
}
