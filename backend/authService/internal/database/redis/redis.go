package redis

import (
	"CRM/go/authService/internal/config"
	"CRM/go/authService/internal/database/postgres"
	"CRM/go/authService/internal/logger"
	"CRM/go/authService/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var (
	rdb    *redis.Client
	pubsub *redis.PubSub
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.GetConfig().Redis.Address,
		Password: config.GetConfig().Redis.Password,
		DB:       0,
	})

	_, err := rdb.ConfigSet(context.Background(), "notify-keyspace-events", "Ex").Result()
	if err != nil {
		log.Fatalf("error setting redis config: %v", err)
	}

	status := rdb.Ping(context.Background())
	if status.Err() != nil {
		log.Fatalf("error connection to redis: %v", status.Err())
	}

	pubsub = rdb.Subscribe(context.Background(), "__keyevent@0__:expired")

	keys, err := rdb.Keys(context.Background(), "*").Result()
	if err != nil {
		panic(err)
	}

	for _, key := range keys {
		fmt.Println("key redis: " + key)
		Del(context.Background(), key)
	}

	go handleExpiredKeys()
}

func PrintKeys() {
	keys, err := rdb.Keys(context.Background(), "*").Result()
	if err != nil {
		panic(err)
	}

	for _, key := range keys {
		fmt.Println("key redis: " + key)
		Del(context.Background(), key)
	}
}

func handleExpiredKeys() {
	for {
		msg, err := ReceiveMessages()
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("error receiving message: %v", err), "message", msg)
			continue
		}
		if msg[:4] == "exp:" {
			var account models.Account
			err = json.Unmarshal([]byte(Get(context.Background(), msg[4:]).Val()), &account)
			if err != nil {
				logger.CreateLog("error", fmt.Sprintf("error unmarshal account: %v", err), "account", account)
			}
			err, _ = postgres.UpdateLastActivityByAccountId(account.Id, account.LastActivity)
			Del(context.Background(), msg[4:])
			if err != nil {
				logger.CreateLog("error", fmt.Sprintf("error deleting key: %v", err), "key", msg[4:])
			}
			logger.CreateLog("info", "key deleted", "key", msg[4:])
		}
		time.Sleep(time.Second)
	}
}

func ReceiveMessages() (string, error) {
	msg, err := pubsub.ReceiveMessage(context.Background())
	if err != nil {
		return "", err
	}
	return msg.Payload, nil
}

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return rdb.Set(ctx, key, value, expiration)
}

func Get(ctx context.Context, key string) *redis.StringCmd {
	return rdb.Get(ctx, key)
}

func Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return rdb.Del(ctx, keys...)
}
