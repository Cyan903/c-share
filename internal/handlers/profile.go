package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Cyan903/c-share/internal/cache"
	"github.com/Cyan903/c-share/internal/database"
	"github.com/Cyan903/c-share/pkg/api"
	"github.com/Cyan903/c-share/pkg/config"
	"github.com/Cyan903/c-share/pkg/log"
	"github.com/Cyan903/c-share/pkg/mail"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func UpdateNickname(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)
	response := api.SimpleResponse{Writer: w}
	nick := r.URL.Query().Get("nickname")

	if api.InvalidNickname(nick) {
		response.BadRequest("Invalid nickname!")
		return
	}

	if err := database.ChangeNickname(id.Issuer, nick); err != nil {
		response.InternalError()
		return
	}

	response.Success("Nickname has been changed!")
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)
	response := api.SimpleResponse{Writer: w}

	oldpass := r.URL.Query().Get("password")
	newpass := r.URL.Query().Get("replacement")

	if api.InvalidPassword(newpass) {
		response.Unauthorized("Bad new password!")
		return
	}

	if err := database.ChangePassword(id.Issuer, oldpass, newpass); err != nil {
		if errors.Is(database.ErrBadPW, err) {
			response.Unauthorized("Invalid password!")
			return
		}

		response.InternalError()
		return
	}
}

// Current email must be verified
func UpdateEmail(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// TODO: User shouldn't be able to spam this
func SendVerification(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)
	response := api.SimpleResponse{Writer: w}

	abt, err := database.About(id.Issuer)
	verifyCode := strings.Replace(uuid.New().String(), "-", "", -1)

	if err != nil {
		response.InternalError()
		return
	}

	if abt.EmailVerified == 1 {
		response.BadRequest("Your email is already verified!")
		return
	}

	if err := cache.SaveVerification(id.Issuer, verifyCode); err != nil {
		response.InternalError()
		return
	}

	m := mail.MailClient{
		To: []string{
			abt.Email,
		},

		From:     config.Data.Mail.User,
		Password: config.Data.Mail.Password,
		Host:     config.Data.Mail.Host,
		Port:     config.Data.Mail.Port,
	}

	// TODO: Improve this
	if err := m.SendMail([]byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: Email Verification\r\n\r\n"+
			"Please verify your email!\n\n%s\r\n",

		config.Data.Mail.User,
		abt.Email, verifyCode,
	))); err != nil {
		log.Error.Println("Could not send email -", err)
		response.InternalError()
		return
	}

	response.Success("Email has been sent!")
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)

	response := api.SimpleResponse{Writer: w}
	code, err := cache.GetVerification(uid.Issuer)

	if err != nil {
		response.InternalError()
		return
	}

	if id != code {
		response.Unauthorized("Invalid code!")
		return
	}

	if err := database.VerifyUserEmail(uid.Issuer); err != nil {
		response.InternalError()
		return
	}

	if err := cache.DeleteEmailVerification(uid.Issuer); err != nil {
		response.InternalError()
		return
	}

	response.Success("Email verified!")
}
