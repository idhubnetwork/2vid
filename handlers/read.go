package handler

import (
	"2vid/logger"
	"2vid/mysql"
	"2vid/redis"
	"errors"
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
	if did != jwt_aud {
		c.JSON(http.StatusForbidden, ActionErr{READ_ERROR})
		return
	}

	jwt_jti, ok := jt.Get("jwt_jti").(string)
	if ok {
		credential, err := getCrentialFromRedis(jwt_iss, jwt_sub, jwt_aud, jwt_jti)
		if err == nil {
			c.JSON(http.StatusOK, credential)
			return
		}
		credential, err = db_mysql.GetCredential(jwt_iss, jwt_sub, jwt_aud, jwt_jti)
		if err != nil {
			c.JSON(http.StatusForbidden, ActionErr{err.Error()})
			return
		}
		c.JSON(http.StatusOK, credential)
		return
	}

	credential, err := getCrentialFromRedis(jwt_iss, jwt_sub, jwt_aud)
	if err == nil {
		c.JSON(http.StatusOK, credential)
		return
	}

	credentials, err := db_mysql.GetCredentials(jwt_iss, jwt_sub, jwt_aud)
	if err != nil {
		c.JSON(http.StatusForbidden, ActionErr{err.Error()})
		return
	}

	c.JSON(http.StatusOK, credentials)
	return
}

func getCrentialFromRedis(args ...string) (*db_mysql.Credential, error) {
	cacheCredential, err := db_redis.GetCacheCredential(args)
	if err != nil {
		logger.Log.Error(err)
		return nil, err
	}
	logger.Log.Debug(cacheCredential)
	if len(cacheCredential.Credential) == 0 {
		return nil, errors.New("no found in redis and go to mysql")
	}
	credential, err := db_mysql.JwtToCredential(cacheCredential.Credential)
	if err != nil {
		logger.Log.Error(err)
		return nil, err
	}
	return credential, nil
}
