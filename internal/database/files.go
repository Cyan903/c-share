package database

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/Cyan903/c-share/pkg/log"
)

type File struct {
	ID          string `json:"id"`
	User        int    `json:"user"`
	FileSize    int64  `json:"file_size"`
	FileType    string `json:"file_type"`
	Permissions int    `json:"permissions"`
	CreatedAt   string `json:"created_at"`
}

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

func GetFile(id string) (File, error) {
	var file File
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	query := Conn.QueryRowContext(c, "SELECT id, user, file_size, file_type, permissions, created_at FROM files WHERE id = ?", id)

	defer cancel()

	if err := query.Scan(
		&file.ID,
		&file.User,
		&file.FileSize,
		&file.FileType,
		&file.Permissions,
		&file.CreatedAt,
	); err != nil {
		log.Error.Println("Error getting file -", err)
		return file, err
	}

	return file, nil
}
