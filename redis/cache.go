package db_redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/idhubnetwork/jsontokens/crypto"
)

// Get a cache credential from redis by jwt_iss, jwt_aud, jwt_sub or jwt_jti,
//   if not exist return nil or if redis throw error return error.
func GetCacheCredential(args []string) (*CacheCredential, error) {
	var data string
	var cacheCredential *CacheCredential

	for _, v := range args {
		data = data + v
	}
	key := string(crypto.SignHash([]byte(data)))

	value, err := redis.Values(DB_redis.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}

	if err := redis.ScanStruct(value, cacheCredential); err != nil {
		return nil, err
	}

	return cacheCredential, nil
}

// Set a cache credential to redis by jwt_iss, jwt_aud, jwt_sub or jwt_jti.
func SetCacheCredential(credential *CacheCredential, args ...string) error {
	var data string

	for _, v := range args {
		data = data + v
	}
	key := string(crypto.SignHash([]byte(data)))

	_, err := DB_redis.Do("HMSET", redis.Args{}.Add(key).AddFlat(credential)...)
	if err != nil {
		return err
	}
	return nil
}

// Delete a pair key-value in redis.
func DelCacheCredential(args ...string) error {
	var data string

	for _, v := range args {
		data = data + v
	}
	key := string(crypto.SignHash([]byte(data)))

	_, err := DB_redis.Do("DEL", key)

	if err != nil {
		return err
	}
	return nil
}
