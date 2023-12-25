package storage

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-spider/config"
	"strconv"
	"time"
)

// 声明一个全局的redisdb变量
var RedisDB *redis.Client

func InitRedisDB() {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     config.ConfData.RedisDB.Host + ":" + strconv.Itoa(config.ConfData.RedisDB.Port),
		Password: config.ConfData.RedisDB.Password, // no password set
		DB:       0,                                // use default DB
	})
	//context.WithTimeout 创建了一个带有 5 秒超时的上下文。如果 Ping 在 5 秒内无法成功，将返回错误
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := RedisDB.Ping(ctx).Result()
	if err != nil {
		fmt.Println("RedisDB.Ping().Result()")
		panic(err)
	}
}

func CloseRedisDB() {
	err := RedisDB.Close()
	if err != nil {
		panic(err)
	}
}
