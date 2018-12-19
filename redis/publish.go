package db_redis

import (
	"encoding/json"
)

// Send message to corresponding channel.
func Publish(channel string, jwt_id int, status int, credential string) error {
	switch channel {
	case "create":
		DB_redis.Do("PUBLISH", channel, credential)
	case "update":
		DB_redis.Do("PUBLISH", channel, jwt_id)
	case "delete":
		DB_redis.Do("PUBLISH", channel, jwt_id)
	case "update_tbd":
		msg, err := getMessage(jwt_id, status, credential)
		if err != nil {
			return err
		}
		DB_redis.Do("PUBLISH", channel, msg)
	case "delete_tbd":
		msg, err := getMessage(jwt_id, status, credential)
		if err != nil {
			return err
		}
		DB_redis.Do("PUBLISH", channel, msg)
	}
	return nil
}

// Redis publish message generator.
func getMessage(jwt_id int, status int, credential string) (string, error) {
	if credential != "" && len(credential) != 0 {
		tmp := CacheCredential{
			status,
			jwt_id,
			credential,
		}
		msg, err := json.Marshal(tmp)
		if err != nil {
			return "", err
		}
		return string(msg), nil
	}
	tmp := struct {
		Status int `json:"status"`
		Jwt_id int `json:"jwt_id"`
	}{
		status,
		jwt_id,
	}
	msg, err := json.Marshal(tmp)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}
