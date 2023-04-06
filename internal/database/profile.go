package database

import (
	"context"

	"github.com/Cyan903/c-share/pkg/log"
)

func ChangeNickname(uid, nick string) error {
	c, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	_, err := Conn.ExecContext(c, "UPDATE users SET nickname = ? WHERE id = ?", nick, uid)

	defer cancel()

	if err != nil {
		log.Info.Println("Could not update user nickname -", err)
		return err
	}

	return nil
}
