package cache

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Cyan903/c-share/pkg/log"
	"github.com/redis/go-redis/v9"
)

type ServerInfo struct {
	Users   int    `json:"users"`
	Storage string `json:"storage"`
	Total   int    `json:"total_files"`
}

func SaveServerStats(users int, storage string, files int) error {
	if err := Conn.Set(
		context.Background(),
		"c:server_info",
		fmt.Sprintf("%d, %s, %d", users, storage, files),
		STATS_EXPIRE,
	).Err(); err != nil {
		log.Error.Println("Could not save server stats to cache -")
		return err
	}

	return nil
}

func GetServerInfo() (ServerInfo, bool, error) {
	var info ServerInfo

	stats, err := Conn.Get(context.Background(), "c:server_info").Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return info, false, nil
		}

		log.Error.Println("Could not get server stats -", err)
		return info, false, err
	}

	ss := strings.Split(stats, ", ")

	users, uerr := strconv.Atoi(ss[0])
	total, terr := strconv.Atoi(ss[2])

	if uerr != nil || terr != nil {
		log.Error.Println("Could not convert server stats -", uerr, terr)
		return info, false, terr
	}

	info.Users = users
	info.Storage = ss[1]
	info.Total = total

	return info, true, nil
}
