package db_redis

import (
	"log"
	"os"

	"github.com/garyburd/redigo/redis"
)

type CacheCredential struct {
	Status     int64  `redis:"status"`
	Jwt_id     int64  `redis:"jwt_id"`
	Credential string `redis:"credential"`
}

var DB_redis redis.Conn

func init() {
	DB_redis, err := redis.DialURL(os.Getenv("REDIS_URL"))
	if err != nil {
		// handle connection error
		log.Fatal(err)
	}

	defer DB_redis.Close()
}
