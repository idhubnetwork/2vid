package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/idhubnetwork/jsontokens"
)

// Unique identify 2vid server in did json token.
const DESTINATION = "url"

// Authentication Error
type AuthErr struct {
	Authentication string `json:"FaliedAuthentication"`
}

// Binding Authentication Token
type Token struct {
	Token string `form:"token" json:"token" xml:"token" binding:"required"`
}

// Gin middleware verify did json token.
//
// Json token authority DID credetntial CRUD action.
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		tmp, err := searchToken(c)

		if err != nil {
			c.JSON(http.StatusForbidden, AuthErr{err.Error()})
			c.Abort()
			return
		}

		jt := jsontokens.NewJsonToken()
		jt.SetToken(tmp)
		if err = jt.Verify(); err != nil {
			c.JSON(http.StatusForbidden, AuthErr{err.Error()})
			c.Abort()
			return
		}
		expiration, ok := jt.Get("expiration").(int)
		if !ok {
			c.JSON(http.StatusForbidden, AuthErr{"jsontoken non expiration"})
			c.Abort()
			return
		}
		if int64(expiration) > time.Now().Unix() {
			c.JSON(http.StatusForbidden, AuthErr{"jsontoken expired"})
			c.Abort()
			return
		}
		destination, ok := jt.Get("destination").(string)
		if !ok {
			c.JSON(http.StatusForbidden, AuthErr{"jsontoken non destination"})
			c.Abort()
			return
		}
		if destination != DESTINATION {
			c.JSON(http.StatusForbidden, AuthErr{"invalid access url destination"})
			c.Abort()
			return
		}
		c.Set("DIDJsonToken", jt)
	}
}

// Json token storage at HTTP Authorization or Field token.
func searchToken(c *gin.Context) (string, error) {
	tmp := c.Request.Header.Get("Authorization")
	var token Token
	if len(tmp) < 14 || tmp[0:13] != "DIDJsonToken " {
		if err := c.Bind(&token); err != nil {
			return "", errors.New("non DID Json Token")
		}
		tmp = token.Token
		if len(tmp) < 14 || string(tmp[0:13]) != "DIDJsonToken " {
			return "", errors.New("invalid DID Json Token")
		}
	}
	return string(tmp[13:]), nil
}
