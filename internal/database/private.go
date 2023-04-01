package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Cyan903/c-share/pkg/log"
)

var PageLen = 10

type FileData struct {
	ID          string `json:"id"`
	FileSize    int64  `json:"file_size"`
	FileType    string `json:"file_type"`
	FileComment string `json:"file_comment"`
	Permissions int    `json:"permissions"`
	CreatedAt   string `json:"created_at"`
}

func FileListing(uid string, page int, perm, fileType, order, sort, comment string) ([]FileData, int, error) {
	var files []FileData
	var pages int
	var fileFilter = "file_type ="
	var searchFilter = "file_comment LIKE"

	perms := map[string]string{
		"any":      "",
		"public":   "AND permissions = 0",
		"private":  "AND permissions = 1",
		"unlisted": "AND permissions = 2",
	}

	orders := map[string]string{
		"any":        "user",
		"size":       "file_size",
		"type":       "file_type",
		"comment":    "file_comment",
		"permission": "permissions",
		"date":       "created_at",
	}

	if fileType == "any" {
		fileFilter = "id !="
	}

	if comment == "" {
		searchFilter = "id !="
	}

	search := fmt.Sprintf(
		`
			SELECT id, file_size, file_type, file_comment, permissions, created_at FROM files
			WHERE user = ? %s AND %s ? AND (%s ?)
			ORDER BY %s %s
			LIMIT ?, %d;
		`, perms[perm], fileFilter, searchFilter, orders[order], sort, PageLen,
	)

	count := fmt.Sprintf(
		`
			SELECT COUNT(1) FROM files
			WHERE user = ? %s AND %s ?
		`, perms[perm], fileFilter,
	)

	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	row, err := Conn.QueryContext(c, search, uid, fileType, "%"+comment+"%", page*PageLen)

	defer cancel()

	if err != nil {
		log.Error.Println("Could not query files -", err)
		return files, 0, err
	}

	for row.Next() {
		var file FileData

		if err := row.Scan(
			&file.ID,
			&file.FileSize,
			&file.FileType,
			&file.FileComment,
			&file.Permissions,
			&file.CreatedAt,
		); err != nil {
			log.Error.Println("Could not query files -", err)
			return files, 0, err
		}

		files = append(files, file)
	}

	if err := Conn.QueryRowContext(c, count, uid, fileType).Scan(&pages); err != nil {
		log.Error.Println("Could not count files -", err)
		return files, 0, err
	}

	return files, pages, nil
}

func GetPrivateFile(id, user string) (File, error) {
	var file File
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	query := Conn.QueryRowContext(c,
		"SELECT id, user, file_size, file_type, file_pass, file_comment, permissions, created_at FROM files WHERE id = ? AND user = ?",
		id, user,
	)

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

	return file, nil
}

func FileInfo(uid, fileID string) (FileData, error) {
	var file FileData
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	query := Conn.QueryRowContext(c,
		"SELECT id, file_size, file_type, file_comment, permissions, created_at FROM files WHERE id = ? AND user = ?",
		fileID, uid,
	)

	defer cancel()

	if err := query.Scan(
		&file.ID,
		&file.FileSize,
		&file.FileType,
		&file.FileComment,
		&file.Permissions,
		&file.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return file, ErrNotFound
		}

		log.Error.Println("Error fetching file info -", err)
		return file, err
	}

	return file, nil
}

func UpdateFileInfo(id, user, password, comment string, permission int) (bool, error) {
	_, err := GetPrivateFile(id, user)

	if err != nil {
		return false, err
	}

	return true, nil
}
