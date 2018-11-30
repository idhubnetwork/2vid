package main

import (
	"2vid/mysql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idhubnetwork/jsontokens"
)

const (
	READ_ERROR = "Only credential audience can read!"
)

func handleCredential(c *gin.Context) {
	jt, ok := c.Get(jt)
	if !ok {
		c.JSON(http.StatusForbidden, "non DID Json Token")
	}

	action, ok := jt.Get("action").(string)
	if !ok {
		c.JSON(http.StatusForbidden, "jsontoken non action")
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
		c.JSON(http.StatusForbidden, "jsontoken invalid action")
	}
}

func readCredential(c *gin.Context, jt *jsontokens.JsonToken) {
	did, ok := jt.Get("did").(string)
	if !ok || len(did) != 32 {
		c.JSON(http.StatusForbidden, "jsontoken invalid or non did")
	}
	jwt_iss, ok := jt.Get("jwt_iss").(string)
	if !ok || len(jwt_iss) != 32 {
		c.JSON(http.StatusForbidden, "jsontoken invalid or non jwt_iss")
	}
	jwt_aud, ok := jt.Get("jwt_aud").(string)
	if !ok || len(jwt_aud) != 32 {
		c.JSON(http.StatusForbidden, "jsontoken invalid or non jwt_aud")
	}
	jwt_sub, ok := jt.Get("jwt_sub").(string)
	if !ok {
		c.JSON(http.StatusForbidden, "jsontoken invalid or non jwt_sub")
	}
	if did != jwt_aud {
		c.JSON(http.StatusForbidden, READ_ERROR)
	}
	jwt_jti, ok := jt.Get("jwt_jti").(string)

	if ok {
		credential, err := db_mysql.GetCredential(jwt_iss, jwt_sub, jwt_aud, jwt_jti)
		if err != nil {
		}
		return
	}
	credentials, err := db_mysql.GetCredentials(jwt_iss, jwt_sub, jwt_aud)
	if err != nil {
	}
	return
}
