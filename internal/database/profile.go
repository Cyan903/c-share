package database

import (
	"context"
	"errors"

	"github.com/Cyan903/c-share/pkg/log"
	"golang.org/x/crypto/bcrypt"
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

func ChangePassword(uid, oldpw, newpw string) error {
	var pw string

	c, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	query := Conn.QueryRowContext(c, "SELECT pw_bcrypt FROM users WHERE id = ?", uid)

	defer cancel()

	if err := query.Scan(&pw); err != nil {
		log.Info.Println("Could not get user's password -", err)
		return err
	}

	// Old password matches current password?
	if err := bcrypt.CompareHashAndPassword([]byte(pw), []byte(oldpw)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrBadPW
		}

		log.Error.Println("Could not compare passwords -", err)
		return err
	}

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(newpw), bcrypt.DefaultCost)

	if err != nil {
		log.Error.Println("Could not hash password -", err)
		return err
	}

	_, err = Conn.ExecContext(c, "UPDATE users SET pw_bcrypt = ? WHERE id = ?", hashedPw, uid)

	if err != nil {
		log.Info.Println("Could not update user password -", err)
		return err
	}

	return nil
}

func VerifyUserEmail(uid string) (error) {
	c, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	_, err := Conn.ExecContext(c, "UPDATE users SET email_verified = 1 WHERE id = ?", uid)

	defer cancel()

	if err != nil {
		log.Info.Println("Could not verify user's email -", err)
		return err
	}

	return nil
}

func ChangeEmail(uid, naddress string) error {
	c, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	_, err := Conn.ExecContext(c, "UPDATE users SET email = ?, email_verified = 0 WHERE id = ?", naddress, uid)

	defer cancel()

	if err != nil {
		log.Error.Println("Could not update user's email -", err)
		return err
	}
	
	return nil
}
