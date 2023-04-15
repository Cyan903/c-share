package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/Cyan903/c-share/pkg/log"
	"github.com/redis/go-redis/v9"
)

func SaveVerification(uid, code string) error {
	if err := Conn.Set(
		context.Background(),
		fmt.Sprintf("c:email_verification:%s", uid),
		code,
		EMAIL_EXPIRE,
	).Err(); err != nil {
		log.Error.Println("Could not save verification code to cache -")
		return err
	}

	return nil
}

func GetVerification(uid string) (string, error) {
	code, err := Conn.Get(context.Background(), fmt.Sprintf("c:email_verification:%s", uid)).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		log.Error.Println("Could not get verification code -", err)
		return "", err
	}

	return code, nil
}

func DeleteEmailVerification(uid string) error {
	if err := Conn.Del(context.Background(), fmt.Sprintf("c:email_verification:%s", uid)).Err(); err != nil {
		log.Error.Println("Could not remove email verification from cache -", err)
		return err
	}

	return nil
}
