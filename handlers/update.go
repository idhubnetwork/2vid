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

	// 0001 1000
	UPDATE_AUDIENCE_OP = 0x18

	UPDATE_ERROR = "Only credential issuer can update!"
)

type JWT struct {
	JsonWebToken string `json:"jwt" binding:"required"`
}

func updateCredential(c *gin.Context, jt *jsontokens.JsonToken) {
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
	if UPDATE_ISSUER_OP&status == 0 {
		if did != jwt_iss {
			c.JSON(http.StatusForbidden, ActionErr{UPDATE_ERROR})
		}
		var jwt JWT
		err := c.ShouldBind(&jwt)
		if err != nil {
			c.JSON(http.StatusForbidden, ActionErr{"invalid or non updated jwt"})
		}
		db_mysql.UpdateCredential_TBD(jwt_id, jwt.JsonWebToken)
		c.JSON(http.StatusOK, ActionSuccess{"jwt update successed but to be determined"})
	}
	if UPDATE_AUDIENCE_OP&status == 0 {
		db_mysql.UpdateCredential(jwt_id)
		c.JSON(http.StatusOK, ActionSuccess{"jwt update successed"})
	}
	c.JSON(http.StatusBadRequest, ActionErr{"invalid update opration"})
}
