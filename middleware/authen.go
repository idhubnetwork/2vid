package middleware

import (
	"strings"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/idhubnetwork/jsontokens"
)

type Token struct {
	token string `form:"token" json:"token" xml:"token" binding:"required"`
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		tmp, err := searchToken(c)
		if err != nil {}
		jt := jsontokens.NewJsonToken()
		jt.SetToken(tmp)
		if err = jt.Verify(), err != nil {}
		c.Set("DIDJsonToken", jt)
	}
}

func searchToken(c *gin.Context) (string, error) {
	tmp := c.Request.Header.Get("Authorization")
	var token Token
	if len(tmp) < 14 || tmp[0, 13] != "DIDJsonToken " {
		if err := c.ShouldBind(&token), err != nil {
			return "", errors.New("non DID Json Token")
		}
		tmp = token.token 
		if len(tmp) < 14 || tmp[0, 13] != "DIDJsonToken " {
			return "", errors.New("invalid DID Json Token")
		}
		return "", errors.New("invalid DID Json Token")
	}
	return string(tmp[13:]), nil
}
