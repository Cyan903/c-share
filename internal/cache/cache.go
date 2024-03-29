package cache

import (
	"context"
	"time"

	"github.com/Cyan903/c-share/pkg/log"
	"github.com/redis/go-redis/v9"
)

const EMAIL_EXPIRE = 2 * time.Hour
const STATS_EXPIRE = 4 * time.Hour

var Conn *redis.Client

func OpenCache(addr, password string, db int) error {
	ctx := context.Background()
	Conn = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	})

	pong, err := Conn.Ping(ctx).Result()

	if err != nil {
		log.Error.Println("Could not ping Redis -", err)
		return err
	}

	log.Info.Println("Connected to cache -", pong)
	return nil
}
