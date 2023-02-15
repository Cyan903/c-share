package database

import (
	"context"
	"time"

	"github.com/Cyan903/c-share/pkg/log"
)

func UploadFile(uid string, size int64, fileType string, permission int) error {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	_, err := Conn.ExecContext(c, "INSERT INTO files VALUES (0, ?, ?, ?, ?, CURRENT_TIMESTAMP)", uid, size, fileType, permission)

	defer cancel()

	if err != nil {
		log.Error.Println("Error inserting file -", err)
		return err
	}

	return nil
}
