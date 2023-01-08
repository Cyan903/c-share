package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/Cyan903/c-share/pkg/log"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"pw_bcrypt"`
}

func EmailUsed(email string) (bool, error) {
	var inUse bool
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	query := Conn.QueryRowContext(c, "SELECT 1 FROM users WHERE email = ?", email)

	defer cancel()

	if err := query.Scan(&inUse); err != nil && err != sql.ErrNoRows {
		log.Error.Printf("Error in EmailUsed | %s\n", err.Error())
		return true, err
	}

	return inUse, nil
}

func Register(nickname, email, password string) (int64, error) {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err != nil {
		log.Error.Println("Could not hash password!", err)
		return 0, err
	}

	user, err := Conn.ExecContext(c, "INSERT INTO users VALUES (0, ?, ?, ?, CURRENT_TIMESTAMP)", nickname, email, hashedPw)
	uid, uerr := user.LastInsertId()

	if err != nil || uerr != nil {
		log.Error.Printf("Could not register user | %s | %s\n", err.Error(), uerr.Error())
		return 0, err
	}

	return uid, nil
}
