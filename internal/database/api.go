package database

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"

	"github.com/Cyan903/c-share/pkg/log"
	"golang.org/x/crypto/bcrypt"
)

type APIToken struct {
	Token     string `json:"token"`
	UserID    string `json:"user"`
	CreatedAt string `json:"created_at"`
}

func GenerateToken(email string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)

	if err != nil {
		log.Error.Println("Could not generate API token", err)
		return "", err
	}

	hasher := md5.New()
	hasher.Write(hash)

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func UserAPIToken(user string) (APIToken, error) {
	var key APIToken

	c, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	query := Conn.QueryRowContext(c, "SELECT user, token, created_at FROM api WHERE user = ?", user)

	defer cancel()

	if err := query.Scan(
		&key.UserID,
		&key.Token,
		&key.CreatedAt,
	); err != nil && err != sql.ErrNoRows {
		log.Error.Println("Could not get API key from user -", err)
		return key, err
	}

	return key, nil
}

func TokenUserData(token string) (APIToken, error) {
	var key APIToken

	c, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	query := Conn.QueryRowContext(c, "SELECT user, token, created_at FROM api WHERE token = ?", token)

	defer cancel()

	if err := query.Scan(
		&key.UserID,
		&key.Token,
		&key.CreatedAt,
	); err != nil && err != sql.ErrNoRows {
		log.Error.Println("Could not get API key from token -", err)
		return key, err
	}

	return key, nil
}

func AddAPIToken(user, token string) error {
	c, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	_, err := Conn.ExecContext(c, "INSERT INTO api VALUES (?, ?, CURRENT_TIMESTAMP)", user, token)

	defer cancel()

	if err != nil {
		log.Error.Println("Error inserting API token -", err)
		return err
	}

	return nil
}

func DeleteAPIToken(user string) error {
	c, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	_, err := Conn.ExecContext(c, "DELETE FROM api WHERE user = ?", user)

	defer cancel()

	if err != nil {
		log.Error.Println("Error removing API token -", err)
		return err
	}

	return nil
}
