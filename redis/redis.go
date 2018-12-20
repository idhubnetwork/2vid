package db_redis

import (
	"context"

	"2vid/config"

	"github.com/garyburd/redigo/redis"
)

type CacheCredential struct {
	Status     int    `redis:"status" json:"status"`
	Jwt_id     int    `redis:"jwt_id" json:"jwt_id"`
	Credential string `redis:"credential" json:"credential"`
}

var DB_redis redis.Conn

func init() {
	var err error
	url := config.V.Redis.Url
	password := config.V.Redis.Password

	DB_redis, err = redis.Dial("tcp", url, redis.DialPassword(password))
	if err != nil {
		// handle connection error
		panic(err)
	}

	ctx, _ := context.WithCancel(context.Background())
	go Subscribe(ctx, url, password, "create", "update", "delete",
		"update_tbd", "delete_tbd")
}
