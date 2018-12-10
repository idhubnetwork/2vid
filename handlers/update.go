package handler

import (
	"2vid/mysql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idhubnetwork/jsontokens"
)

const (
	// 0000 0010
	UPDATE_ISSUER_OP = 0x02

	// 1110 0111
	UPDATE_ISSUER_OP_TBD = 0xe7

	// 0001 1000
	UPDATE_AUDIENCE_OP = 0x18

	// 0001 1010 & 0011 0101 = 0001 0000 0x10
	IF_CAN_NOT_UPDATE = 0x1a

	// 0001 0000
	CAN_NOT_UPDATE = 0x10

	UPDATE_ERROR = "Only credential issuer can update!"

	UPDATE_NEED_ERROR = "Need credential audience agree update!"
)

type JWT struct {
	JsonWebToken string `json:"jwt" binding:"required"`
}

// Update credential, 2 cases:
//
// both iss and aud can not update
// iss update but need aud agree
func updateCredential(c *gin.Context, jt *jsontokens.JsonToken) {
	did, ok := jt.Get("did").(string)
	if !ok || len(did) != 52 {
		c.JSON(http.StatusForbidden, ActionErr{"jsontoken invalid or non did"})
	}
	jwt_iss, ok := jt.Get("jwt_iss").(string)
	if !ok || len(jwt_iss) != 52 {
		c.JSON(http.StatusForbidden, ActionErr{"jsontoken invalid or non jwt_iss"})
	}
	jwt_aud, ok := jt.Get("jwt_aud").(string)
	if !ok || len(jwt_aud) != 52 {
		c.JSON(http.StatusForbidden, ActionErr{"jsontoken invalid or non jwt_aud"})
	}
	jwt_sub, ok := jt.Get("jwt_sub").(string)
	if !ok {
		c.JSON(http.StatusForbidden, ActionErr{"jsontoken invalid or non jwt_sub"})
	}

	var (
		jwt_id int
		status int
		err    error
	)
	jwt_jti, ok := jt.Get("jwt_jti").(string)
	if !ok {
		jwt_id, status, err = db_mysql.GetStatus(jwt_iss, jwt_sub, jwt_aud)
		if err != nil {
			c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		}
	} else {
		jwt_id, status, err = db_mysql.GetStatus(jwt_iss, jwt_sub, jwt_aud, jwt_jti)
		if err != nil {
			c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		}
	}

	if IF_CAN_NOT_UPDATE&status == CAN_NOT_UPDATE {
		c.JSON(http.StatusForbidden, ActionErr{"This credential can't update"})
	}

	if UPDATE_ISSUER_OP&status == 0 {
		if did != jwt_iss {
			c.JSON(http.StatusForbidden, ActionErr{UPDATE_ERROR})
		}
		tmp, ok := c.Get("credential")
		jwt, ok := tmp.(string)
		if !ok {
			c.JSON(http.StatusForbidden, ActionErr{"invalid or non updated jwt"})
		}
		credential, err := db_mysql.VerifyWritedData(did, jwt)
		if err != nil {
			c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		}
		status = status & UPDATE_ISSUER_OP_TBD
		db_mysql.UpdateCredential_TBD(jwt_id, status, credential)
		c.JSON(http.StatusOK, ActionSuccess{"jwt update successed but to be determined"})
	}
	if UPDATE_AUDIENCE_OP&status == 0 {
		if did != jwt_aud {
			c.JSON(http.StatusForbidden, ActionErr{UPDATE_NEED_ERROR})
		}
		db_mysql.UpdateCredential(jwt_id)
		c.JSON(http.StatusOK, ActionSuccess{"credential update successed"})
	}
	c.JSON(http.StatusBadRequest, ActionErr{"invalid update opration"})
}
