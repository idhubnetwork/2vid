package db_redis

import (
	"log"
	"os"

	"github.com/garyburd/redigo/redis"
)

type CacheCredential struct {
	Status     int    `redis:"status" json:"status"`
	Jwt_id     int    `redis:"jwt_id" json:"jwt_id"`
	Credential string `redis:"credential" json:"credential"`
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
