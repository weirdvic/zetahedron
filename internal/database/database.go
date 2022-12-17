package database

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func CreateClient(dbNo int) *redis.Client {
	DB_ADDR := os.Getenv("DB_ADDR")
	// Try default value in case DB_ADDR was not provided
	if DB_ADDR == "" {
		DB_ADDR = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASS"),
		DB:       dbNo,
	})
	return rdb
}
