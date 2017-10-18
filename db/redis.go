package db

import (
    "github.com/go-redis/redis"

    "CRBot/config"
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

func addJob(key, value string) (error){
	err := GetInstance().RPush(key, value).Err()
	return err
}

func popJob(key string) (string){
	return GetInstance().LPop(key).Val()
}

func AddPlayerStatsJob(playertag, chatid string) (error){
	value := playertag + ":" + chatid
	return addJob("stats", value)
}

func RemovePlayerStatsJob() (string){
	return popJob("stats")
}
