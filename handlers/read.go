package handler

import (
	"2vid/mysql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idhubnetwork/jsontokens"
)

// Return unique credential or credentials array with json form.
//
// Params {iss, aud, sub, jti} identify a unique credential.
// Params {iss, aud, sub} return a credential array OR maybe should identify
//   a unique credential to be determine.
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
		c.JSON(http.StatusOK, credential)
	}
	credentials, err := db_mysql.GetCredentials(jwt_iss, jwt_sub, jwt_aud)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
	}
	c.JSON(http.StatusOK, credentials)
}
