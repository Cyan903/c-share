package database

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/Cyan903/c-share/pkg/log"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func IDUsed(id string) (bool, error) {
	var inUse bool
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	query := Conn.QueryRowContext(c, "SELECT 1 FROM files WHERE id = ?", id)

	defer cancel()

	if err := query.Scan(&inUse); err != nil && err != sql.ErrNoRows {
		log.Error.Printf("Error in IDUsed | %s\n", err.Error())
		return true, err
	}

	return inUse, nil
}

func RandomID() (string, error) {
	for {
		brand := make([]byte, 10+2)
		rand.Read(brand)

		code := fmt.Sprintf("%x", brand)[2 : 10+2]
		check, err := IDUsed(code)

		if err != nil {
			log.Error.Printf("Error in RandomID | %s\n", err.Error())
			return "", err
		}

		if check {
			log.Warning.Println(code, "was in use. Creating a new one...")
		} else {
			return code, nil
		}
	}
}

func UploadFile(rid, uid string, size int64, fileType string, permission int) error {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	_, err := Conn.ExecContext(c, "INSERT INTO files VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)", rid, uid, size, fileType, permission)

	defer cancel()

	if err != nil {
		log.Error.Println("Error inserting file -", err)
		return err
	}

	return nil
}
