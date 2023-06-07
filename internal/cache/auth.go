package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/Cyan903/c-share/pkg/log"
	"github.com/redis/go-redis/v9"
)

func SaveResetToken(email, token string) error {
	if err := Conn.Set(
		context.Background(),
		fmt.Sprintf("c:password_reset:%s", token),
		email,
		EMAIL_EXPIRE,
	).Err(); err != nil {
		log.Error.Println("Could not save reset token to cache -")
		return err
	}

	return nil
}

func GetResetToken(token string) (string, error) {
	token, err := Conn.Get(context.Background(), fmt.Sprintf("c:password_reset:%s", token)).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		log.Error.Println("Could not get reset token -", err)
		return "", err
	}

	return token, nil
}

func DeleteResetToken(token string) error {
	if err := Conn.Del(context.Background(), fmt.Sprintf("c:password_reset:%s", token)).Err(); err != nil {
		log.Error.Println("Could not remove reset token from cache -", err)
		return err
	}

	return nil
}
