package initializers

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func ConnectRedis() {
	opt, err := redis.ParseURL(os.Getenv("REDIS_CONN_STRING"))
	if err != nil {
		panic("Fail to connect to redis")
	}
	RDB = redis.NewClient(opt)
	log.Print("Connected to Redis")
}
