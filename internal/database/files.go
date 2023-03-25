package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Cyan903/c-share/pkg/log"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
)

type File struct {
	ID          string `json:"id"`
	User        int    `json:"user"`
	FileSize    int64  `json:"file_size"`
	FileType    string `json:"file_type"`
	FilePass    string `json:"file_href"`
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
		log.Error.Println("Could not check ID -", err)
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
			log.Error.Println("Error in RandomID -", err)
			return "", err
		}

		if check {
			log.Warning.Println(code, "was in use. Creating a new one...")
		} else {
			return code, nil
		}
	}
}

func UploadFile(rid, uid string, size int64, fileType, filePass string, permission int) error {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(filePass), bcrypt.DefaultCost)
	if err != nil {
		log.Error.Println("Could not hash password!")
		return err
	}

	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	_, err = Conn.ExecContext(c, "INSERT INTO files VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)", rid, uid, size, fileType, hashedPw, permission)

	defer cancel()

	if err != nil {
		log.Error.Println("Error inserting file -", err)
		return err
	}

	return nil
}

func GetFile(id, pass string) (File, error) {
	var file File
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	query := Conn.QueryRowContext(c, "SELECT id, user, file_size, file_type, file_pass, permissions, created_at FROM files WHERE id = ?", id)

	defer cancel()

	if err := query.Scan(
		&file.ID,
		&file.User,
		&file.FileSize,
		&file.FileType,
		&file.FilePass,
		&file.Permissions,
		&file.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return file, ErrNotFound
		}

		log.Error.Println("Error fetching file -", err)
		return file, err
	}

	if file.Permissions == 2 {
		if err := bcrypt.CompareHashAndPassword([]byte(file.FilePass), []byte(pass)); err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				return file, ErrBadPW
			}

			log.Error.Println("Could not compare passwords -", err)
			return file, err
		}
	}

	return file, nil
}

func OwnFiles(id []string, uid int) ([]string, error) {
	var files []File
	var dbIDs []string
	var ids string
	var notOwned []string

	args := make([]interface{}, len(id))

	for i, iid := range id {
		args[i] = iid
	}

	if len(id) > 1 {
		ids = strings.Repeat("?, ", len(id)-1) + "?"
	} else {
		ids = "?"
	}

	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	row, err := Conn.QueryContext(
		c, fmt.Sprintf("SELECT id, user, file_size, file_type, file_pass, permissions, created_at FROM files WHERE id IN (%s)", ids),
		args...,
	)

	defer cancel()

	if err != nil {
		log.Error.Println("Error checking ownership -", err)
		return notOwned, err
	}

	for row.Next() {
		var file File

		if err := row.Scan(
			&file.ID,
			&file.User,
			&file.FileSize,
			&file.FileType,
			&file.FilePass,
			&file.Permissions,
			&file.CreatedAt,
		); err != nil {
			log.Error.Println("Error scaning OwnFiles -", err)
			return notOwned, err
		}

		files = append(files, file)
		dbIDs = append(dbIDs, file.ID)
	}

	// Does file exist on DB?
	for _, f := range id {
		if !slices.Contains(dbIDs, f) {
			log.Warning.Println(f, "file does not exist in database!")
			notOwned = append(notOwned, f)
		}
	}

	// Does user own files?
	for _, f := range files {
		if f.User != uid {
			notOwned = append(notOwned, f.ID)
		}
	}

	return notOwned, nil
}

func DeleteFiles(uid string, files []string) error {
	var purgeList string

	args := make([]interface{}, len(files)+1)
	args[len(args)-1] = uid

	for i, iid := range files {
		args[i] = iid
	}

	if len(files) > 1 {
		purgeList = strings.Repeat("?, ", len(files)-1) + "?"
	} else {
		purgeList = "?"
	}

	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	_, err := Conn.ExecContext(
		c, fmt.Sprintf("DELETE FROM files WHERE id IN (%s) AND user = ?", purgeList),
		args...,
	)

	defer cancel()

	if err != nil {
		log.Error.Println("Error removing file -", err)
		return err
	}

	return nil
}
