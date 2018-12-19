package handler

import (
	"2vid/logger"
	"2vid/mysql"
	"2vid/redis"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idhubnetwork/jsontokens"
)

// Delete credential, 5 cases:
//
// both iss and aud can not delete.
// iss delete need aud agree
// aud delete need iss agree
// iss delete directly
// aud delete directly
func deleteCredential(c *gin.Context, jt *jsontokens.JsonToken) {
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
	cacheCredential, err := db_redis.GetCacheCredential([]string{jwt_iss, jwt_sub, jwt_aud})
	if cacheCredential == nil || err != nil {
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
	}

	if IF_CAN_NOT_DELETE&status == CAN_NOT_DELETE {
		c.JSON(http.StatusForbidden, ActionErr{"This credential can't delete"})
	}

	if DELETE_ISSUER_OP&status == 0 {
		if did != jwt_iss {
			c.JSON(http.StatusForbidden, ActionErr{DELETE_ISSUER_ERROR})
		}
		err = db_redis.Publish("delete", jwt_id, 0, "")
		if err != nil {
			logger.Log.Error(err)
			c.JSON(http.StatusForbidden, ActionErr{err.Error()})
			return
		}

		c.JSON(http.StatusOK, ActionSuccess{"credential delete successed"})
	}

	if DELETE_AUDIENCE_OP&status == 0 {
		if did != jwt_aud {
			c.JSON(http.StatusForbidden, ActionErr{DELETE_AUDIENCE_ERROR})
		}
		db_mysql.DeleteCredential(jwt_id)
		c.JSON(http.StatusOK, ActionSuccess{"credential delete successed"})
	}

	if did == jwt_iss {
		status = status & DELETE_ISSUER_OP_TBD
		err = db_redis.Publish("delete_tbd", jwt_id, status, "")
		if err != nil {
			logger.Log.Error(err)
			c.JSON(http.StatusForbidden, ActionErr{err.Error()})
			return
		}

		c.JSON(http.StatusOK, ActionSuccess{"credential delete successed but to be determined"})
	}

	if did == jwt_aud {
		status = status & DELETE_AUDIENCE_OP_TBD
		err = db_redis.Publish("delete_tbd", jwt_id, status, "")
		if err != nil {
			logger.Log.Error(err)
			c.JSON(http.StatusForbidden, ActionErr{err.Error()})
			return
		}

		c.JSON(http.StatusOK, ActionSuccess{"credential delete successed but to be determined"})
	}

	c.JSON(http.StatusBadRequest, ActionErr{"invalid delete opration"})
}
