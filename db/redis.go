package dnscache

import (
	"time"

    "github.com/go-redis/redis"

    "GoHole/config"
)

var instance *redis.Client = nil

func GetInstance() *redis.Client {
    if instance == nil {
    	host := config.GetInstance().RedisDB.Host
    	port := config.GetInstance().RedisDB.Port
    	addr := host + ":" + port
        instance = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: config.GetInstance().RedisDB.Pass,
			DB:       0,  // use default DB
		})

		_, err := instance.Ping().Result()
		if err != nil {
			panic(err)
		}
    }

    return instance
}

