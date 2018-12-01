package handler

import (
	"2vid/mysql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idhubnetwork/jsontokens"
)

const (
	READ_ERROR = "Only credential audience can read!"
)

type ActionErr struct {
	ActionError string `json:"FaliedAction"`
}

type ActionSuccess struct {
	Action string `json:"Action"`
}

func handleCredential(c *gin.Context) {
	tmp, ok := c.Get("jt")
	if !ok {
		c.JSON(http.StatusForbidden, ActionErr{"non DID Json Token"})
	}
	jt := tmp.(*jsontokens.JsonToken)

	action, ok := jt.Get("action").(string)
	if !ok {
		c.JSON(http.StatusForbidden, ActionErr{"jsontoken non action"})
	}

	switch action {
	case "READ":
		readCredential(c, jt)
	case "CREATE":
		createCredential(c, jt)
	case "UPDATE":
		updateCredential(c, jt)
	case "DELETE":
		deleteCredential(c, jt)
	default:
		c.JSON(http.StatusForbidden, ActionErr{"jsontoken invalid action"})
	}
}

func readCredential(c *gin.Context, jt *jsontokens.JsonToken) {
	did, ok := jt.Get("did").(string)
	if !ok || len(did) != 32 {
		c.JSON(http.StatusForbidden, ActionErr{"jsontoken invalid or non did"})
	}
	jwt_iss, ok := jt.Get("jwt_iss").(string)
	if !ok || len(jwt_iss) != 32 {
		c.JSON(http.StatusForbidden, ActionErr{"jsontoken invalid or non jwt_iss"})
	}
	jwt_aud, ok := jt.Get("jwt_aud").(string)
	if !ok || len(jwt_aud) != 32 {
		c.JSON(http.StatusForbidden, ActionErr{"jsontoken invalid or non jwt_aud"})
	}
	jwt_sub, ok := jt.Get("jwt_sub").(string)
	if !ok {
		c.JSON(http.StatusForbidden, ActionErr{"jsontoken invalid or non jwt_sub"})
	}
	if did != jwt_aud {
		c.JSON(http.StatusForbidden, ActionErr{READ_ERROR})
	}
	jwt_jti, ok := jt.Get("jwt_jti").(string)

	if ok {
		credential, err := db_mysql.GetCredential(jwt_iss, jwt_sub, jwt_aud, jwt_jti)
		if err != nil {
			c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		}
		c.JSON(http.StatusForbidden, credential)
	}
	credentials, err := db_mysql.GetCredentials(jwt_iss, jwt_sub, jwt_aud)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
	}
	c.JSON(http.StatusForbidden, credentials)
}
