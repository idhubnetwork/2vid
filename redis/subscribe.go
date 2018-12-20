package db_redis

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"2vid/logger"
	"2vid/mysql"

	"github.com/garyburd/redigo/redis"
)

// Subscribe listens for messages on Redis pubsub channels.
func Subscribe(ctx context.Context, redisServerAddr string, password string,
	channels ...string) error {
	// A ping is set to the server with this period to test for the health of
	// the connection and server.
	const healthCheckPeriod = time.Minute

	c, err := redis.DialURL(redisServerAddr, redis.DialPassword(password),
		// Read timeout on server should be greater than ping period.
		redis.DialReadTimeout(healthCheckPeriod+10*time.Second),
		redis.DialWriteTimeout(10*time.Second))
	if err != nil {
		logger.Log.Error(err)
		return err
	}
	defer c.Close()

	psc := redis.PubSubConn{Conn: c}

	if err := psc.Subscribe(redis.Args{}.AddFlat(channels)...); err != nil {
		logger.Log.Error(err)
		return err
	}

	done := make(chan error, 1)

	// Start a goroutine to receive notifications from the server.
	go func() {
		for {
			switch n := psc.Receive().(type) {
			case error:
				logger.Log.Error(n)
				done <- n
				return
			case redis.Message:
				if err := consume(n.Channel, n.Data); err != nil {
					logger.Log.Error(err)
					done <- err
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(healthCheckPeriod)
	defer ticker.Stop()
loop:
	for err == nil {
		select {
		case <-ticker.C:
			// Send ping to test health of connection and server. If
			// corresponding pong is not received, then receive on the
			// connection will timeout and the receive goroutine will exit.
			if err = psc.Ping(""); err != nil {
				logger.Log.Error(err)
				break loop
			}
		case <-ctx.Done():
			break loop
		case err := <-done:
			// Return error from the receive goroutine.
			return err
		}
	}
	return errors.New("subscribe end")
}

func consume(channel string, data []byte) error {
	switch channel {
	case "create":
		credential, err := db_mysql.JwtToCredential(string(data))
		if err != nil {
			return err
		}
		id, err := db_mysql.CreateCredential(credential)
		if err != nil {
			return err
		}
		cacheCredential := CacheCredential{
			credential.Status,
			id,
			credential.Credential,
		}
		err = SetCacheCredential(&cacheCredential, credential.Iss,
			credential.Sub, credential.Aud)
		if err != nil {
			return err
		}
	case "update":
		tmp, err := strconv.Atoi(string(data))
		if err != nil {
			return err
		}
		err = db_mysql.UpdateCredential(tmp)
		if err != nil {
			return err
		}
		credential, err := db_mysql.GetKeyById(tmp)
		if err != nil {
			return err
		}
		cacheCredential := CacheCredential{
			credential.Status,
			tmp,
			credential.Credential,
		}
		err = SetCacheCredential(&cacheCredential, credential.Iss,
			credential.Sub, credential.Aud)
		if err != nil {
			return err
		}
	case "delete":
		tmp, err := strconv.Atoi(string(data))
		if err != nil {
			return err
		}
		credential, err := db_mysql.GetKeyById(tmp)
		if err != nil {
			return err
		}
		err = DelCacheCredential(credential.Iss, credential.Sub, credential.Aud)
		if err != nil {
			return err
		}
		err = db_mysql.DeleteCredential(tmp)
		if err != nil {
			return err
		}
	case "update_tbd":
		tmp := new(CacheCredential)
		err := json.Unmarshal(data, tmp)
		if err != nil {
			return err
		}
		credential, err := db_mysql.JwtToCredential(tmp.Credential)
		if err != nil {
			return err
		}
		err = db_mysql.UpdateCredential_TBD(tmp.Jwt_id, tmp.Status, credential)
		if err != nil {
			return err
		}
		cacheCredential := CacheCredential{
			tmp.Status,
			tmp.Jwt_id,
			credential.Credential,
		}
		err = SetCacheCredential(&cacheCredential, credential.Iss,
			credential.Sub, credential.Aud)
		if err != nil {
			return err
		}
	case "delete_tbd":
		tmp := new(CacheCredential)
		err := json.Unmarshal(data, tmp)
		if err != nil {
			return err
		}
		db_mysql.DeleteCredential_TBD(tmp.Jwt_id, tmp.Status)
		if err != nil {
			return err
		}
		credential, err := db_mysql.GetKeyById(tmp.Jwt_id)
		if err != nil {
			return err
		}
		cacheCredential := CacheCredential{
			credential.Status,
			tmp.Jwt_id,
			credential.Credential,
		}
		err = SetCacheCredential(&cacheCredential, credential.Iss,
			credential.Sub, credential.Aud)
		if err != nil {
			return err
		}
	}
	return nil
}
