package database

import (
	"context"
	"database/sql"
	"errors"
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

type Info struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
	Register string `json:"created_at"`
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

func Login(email, password string) (Users, error) {
	var usr Users

	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	res := Conn.QueryRowContext(c, "SELECT id, nickname, email, pw_bcrypt FROM users WHERE email = ?", email)
	err := res.Scan(
		&usr.ID,
		&usr.Nickname,
		&usr.Email,
		&usr.Password,
	)

	defer cancel()

	if err == sql.ErrNoRows {
		log.Error.Println("User does not exist | ", email)
		return usr, err
	} else if err != nil {
		log.Error.Println("Could not fetch user info |", err)
		return usr, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return usr, err
		}

		log.Error.Println("Could not compare password |", err)
		return usr, err
	}

	return usr, nil
}

func About(uid string) (Info, error) {
	var abt Info
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	res := Conn.QueryRowContext(c, "SELECT id, nickname, created_at FROM users WHERE id = ?", uid)

	defer cancel()

	if err := res.Scan(
		&abt.ID,
		&abt.Nickname,
		&abt.Register,
	); err != nil {
		log.Error.Println("Could not fetch about -", err)
		return abt, err
	}

	return abt, nil
}
