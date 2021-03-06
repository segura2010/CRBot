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

func SubscribeToPlayerStats() (*redis.PubSub){
	return GetInstance().PSubscribe("stats")
}

func PublishToPlayerStats(msg string) (error){
	return GetInstance().Publish("stats", msg).Err()
}

func AddPlayerTag(playertag, telegramid string) (error){
	return GetInstance().Set(telegramid, playertag, 0).Err()
}

func DeletePlayerTag(telegramid string) (error){
	return GetInstance().Del(telegramid).Err()
}

func GetPlayerTag(telegramid string) (string){
	return GetInstance().Get(telegramid).Val()
}
