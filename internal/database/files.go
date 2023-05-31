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
	FilePass    string `json:"file_pass"`
	FileComment string `json:"file_comment"`
	Permissions int    `json:"permissions"`
	CreatedAt   string `json:"created_at"`
}

type StorageInfo struct {
	Users     int    `json:"users"`
	Storage   string `json:"storage"`
	FileCount struct {
		Total    int `json:"total"`
		Public   int `json:"public"`
		Private  int `json:"private"`
		Unlisted int `json:"unlisted"`
	} `json:"files"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func IDUsed(id string) (bool, error) {
	var inUse bool
	c, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
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

func UploadFile(rid, uid string, size int64, fileType, filePass, comment string, permission int) error {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(filePass), bcrypt.DefaultCost)
	if err != nil {
		log.Error.Println("Could not hash password!")
		return err
	}

	c, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	_, err = Conn.ExecContext(c, "INSERT INTO files VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)", rid, uid, size, fileType, hashedPw, comment, permission)

	defer cancel()

	if err != nil {
		log.Error.Println("Error inserting file -", err)
		return err
	}

	return nil
}

func GetFile(id, pass string) (File, error) {
	var file File
	c, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	query := Conn.QueryRowContext(c, "SELECT id, user, file_size, file_type, file_pass, file_comment, permissions, created_at FROM files WHERE id = ?", id)

	defer cancel()

	if err := query.Scan(
		&file.ID,
		&file.User,
		&file.FileSize,
		&file.FileType,
		&file.FilePass,
		&file.FileComment,
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
	var notOwned []string

	ids := "?"
	args := make([]interface{}, len(id))

	for i, iid := range id {
		args[i] = iid
	}

	if len(id) > 1 {
		ids = strings.Repeat("?, ", len(id)-1) + "?"
	}

	c, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	row, err := Conn.QueryContext(
		c, fmt.Sprintf("SELECT id, user, file_size, file_type, file_pass, file_comment, permissions, created_at FROM files WHERE id IN (%s)", ids),
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
			&file.FileComment,
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
	purgeList := "?"
	args := make([]interface{}, len(files)+1)

	args[len(args)-1] = uid

	for i, iid := range files {
		args[i] = iid
	}

	if len(files) > 1 {
		purgeList = strings.Repeat("?, ", len(files)-1) + "?"
	}

	c, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
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

// TODO: check if user has enough storage left
func UpdateStorage(uid string) (int64, error) {
	var storage int64

	c, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	query := Conn.QueryRowContext(c, "SELECT COALESCE(SUM(file_size), 0) FROM files WHERE user = ?", uid)

	defer cancel()

	if err := query.Scan(&storage); err != nil {
		log.Error.Println("Error calculating storage -", err)
		return storage, err
	}

	if _, err := Conn.ExecContext(c, "UPDATE users SET used_storage = ? WHERE id = ?", storage, uid); err != nil {
		log.Error.Println("Error updating user's storage value -", err)
		return storage, err
	}

	return storage, nil
}

func ServerStorageInfo() (StorageInfo, error) {
	var sinfo StorageInfo

	c, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	query := Conn.QueryRowContext(c, `
		SELECT
			COUNT(1) AS users,
			SUM(used_storage) AS storage,
			(SELECT COUNT(1) FROM files) AS total_files,
			(SELECT COUNT(1) FROM files WHERE permissions = 0) AS pub,
			(SELECT COUNT(1) FROM files WHERE permissions = 1) AS priv,
			(SELECT COUNT(1) FROM files WHERE permissions = 2) AS unlist
		FROM users;
	`)

	defer cancel()

	if err := query.Scan(
		&sinfo.Users,
		&sinfo.Storage,
		&sinfo.FileCount.Total,
		&sinfo.FileCount.Public,
		&sinfo.FileCount.Private,
		&sinfo.FileCount.Unlisted,
	); err != nil {
		log.Error.Println("Error fetching storage info -", err)
		return sinfo, err
	}

	return sinfo, nil
}
