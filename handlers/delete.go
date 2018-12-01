package handler

import (
	"2vid/mysql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idhubnetwork/jsontokens"
)

const (
	// 0000 0001
	DELETE_ISSUER_OP = 0x01

	// 0000 0100
	DELETE_AUDIENCE_OP = 0x04

	// 1101 1011
	DELETE_ISSUER_OP_TBD = 0xdb

	// 1101 1110
	DELETE_AUDIENCE_OP_TBD = 0xde

	DELETE_ISSUER_ERROR = "Credential issuer can delete but no authorization!"

	UPDATE_AUDIENCE_ERROR = "Credential audience can delete but no authorization!"
)

func DeleteCredential(c *gin.Context, jt *jsontokens.JsonToken) {
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

	var (
		jwt_id int
		status int
	)
	jwt_jti, ok := jt.Get("jwt_jti").(string)
	if !ok {
		jwt_id, status, err := db_mysql.GetStatus(jwt_iss, jwt_sub, jwt_aud)
		if err != nil {
			c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		}
	} else {
		jwt_id, status, err := db_mysql.GetStatus(jwt_iss, jwt_sub, jwt_aud, jwt_jti)
		if err != nil {
			c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		}
	}

	if DELETE_ISSUER_OP&status == 0 {
		if did != jwt_iss {
			c.JSON(http.StatusForbidden, ActionErr{UPDATE_ISSUER_ERROR})
		}
		db_mysql.DeleteCredential(jwt_id)
		c.JSON(http.StatusOK, ActionSuccess{"jwt delete successed"})
	}

	if DELETE_AUDIENCE_OP&status == 0 {
		if did != jwt_aud {
			c.JSON(http.StatusForbidden, ActionErr{UPDATE_AUDIENCE_ERROR})
		}
		db_mysql.DeleteCredential(jwt_id)
		c.JSON(http.StatusOK, ActionSuccess{"jwt delete successed"})
	}

	if did == jwt_iss {
		status = status & DELETE_ISSUER_OP_TBD
		db_mysql.DeleteCredential_TBD(jwt_id, status)
	}

	if did == jwt_aud {
		status = status & DELETE_AUDIENCE_OP_TBD
		db_mysql.DeleteCredential_TBD(jwt_id, status)
	}

	c.JSON(http.StatusBadRequest, ActionErr{"invalid delete opration"})
}
