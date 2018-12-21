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
		return
	}
	jwt_iss, ok := jt.Get("jwt_iss").(string)
	if !ok || len(jwt_iss) != 52 {
		c.JSON(http.StatusForbidden, ActionErr{"jsontoken invalid or non jwt_iss"})
		return
	}
	jwt_aud, ok := jt.Get("jwt_aud").(string)
	if !ok || len(jwt_aud) != 52 {
		c.JSON(http.StatusForbidden, ActionErr{"jsontoken invalid or non jwt_aud"})
		return
	}
	jwt_sub, ok := jt.Get("jwt_sub").(string)
	if !ok {
		c.JSON(http.StatusForbidden, ActionErr{"jsontoken invalid or non jwt_sub"})
		return
	}

	var (
		jwt_id int
		status int
		err    error
	)
	cacheCredential, err := db_redis.GetCacheCredential([]string{jwt_iss, jwt_sub, jwt_aud})
	if err != nil {
		logger.Log.Error(err)
	}

	logger.Log.Debug(cacheCredential)
	if cacheCredential.Jwt_id != 0 && cacheCredential.Status != 0 {
		jwt_id = cacheCredential.Jwt_id
		status = cacheCredential.Status
	} else {
		logger.Log.Debug(cacheCredential)
		jwt_jti, ok := jt.Get("jwt_jti").(string)
		if !ok {
			jwt_id, status, err = db_mysql.GetStatus(jwt_iss, jwt_sub, jwt_aud)
			if err != nil {
				c.JSON(http.StatusForbidden, ActionErr{err.Error()})
				return
			}
		} else {
			jwt_id, status, err = db_mysql.GetStatus(jwt_iss, jwt_sub, jwt_aud, jwt_jti)
			if err != nil {
				c.JSON(http.StatusForbidden, ActionErr{err.Error()})
				return
			}
		}
	}
	logger.Log.Debug(status)

	if IF_CAN_NOT_DELETE&status == CAN_NOT_DELETE {
		c.JSON(http.StatusForbidden, ActionErr{"This credential can't delete"})
		return
	}

	if DELETE_ISSUER_OP&status == 0 {
		if did != jwt_iss {
			c.JSON(http.StatusForbidden, ActionErr{DELETE_ISSUER_ERROR})
			return
		}
		err = db_redis.Publish("delete", jwt_id, 0, "")
		if err != nil {
			logger.Log.Error(err)
			c.JSON(http.StatusForbidden, ActionErr{err.Error()})
			return
		}

		c.JSON(http.StatusOK, ActionSuccess{"credential delete successed"})
		return
	}

	if DELETE_AUDIENCE_OP&status == 0 {
		if did != jwt_aud {
			c.JSON(http.StatusForbidden, ActionErr{DELETE_AUDIENCE_ERROR})
			return
		}
		db_mysql.DeleteCredential(jwt_id)
		c.JSON(http.StatusOK, ActionSuccess{"credential delete successed"})
		return
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
		return
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
		return
	}

	c.JSON(http.StatusBadRequest, ActionErr{"invalid delete opration"})
	return
}
