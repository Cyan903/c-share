package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Cyan903/c-share/internal/cache"
	"github.com/Cyan903/c-share/internal/database"
	"github.com/Cyan903/c-share/pkg/api"
	"github.com/Cyan903/c-share/pkg/auth"
	"github.com/Cyan903/c-share/pkg/config"
	"github.com/Cyan903/c-share/pkg/log"
	"github.com/Cyan903/c-share/pkg/mail"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type account struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type user struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var acc account

	accDecoder := json.NewDecoder(r.Body)
	response := api.SimpleResponse{Writer: w}
	tokenResponse := api.AdvancedResponse{Writer: w}

	if err := accDecoder.Decode(&acc); err != nil {
		response.BadRequest("Invalid JSON!")
		return
	}

	// Check nickname
	if api.InvalidNickname(acc.Nickname) {
		response.BadRequest("Invalid nickname!")
		return
	}

	// Check password
	if api.InvalidPassword(acc.Password) {
		response.BadRequest("Invalid password!")
		return
	}

	// Check email
	if api.InvalidEmail(acc.Email) {
		response.BadRequest("Invalid email!")
		return
	}

	inUse, err := database.EmailUsed(acc.Email)

	if err != nil {
		response.InternalError()
		return
	}

	if inUse {
		response.Conflict("Email is in use!")
		return
	}

	// Register request
	uid, err := database.Register(acc.Nickname, acc.Email, acc.Password)

	if err != nil {
		response.InternalError()
		return
	}

	// Login
	token, err := auth.Create(strconv.Itoa(int(uid)))

	if err != nil {
		response.InternalError()
		return
	}

	tokenResponse.Code = http.StatusOK
	tokenResponse.Data = token
	tokenResponse.Count = 1

	tokenResponse.JSON()
}

func Login(w http.ResponseWriter, r *http.Request) {
	var usr user

	usrDecoder := json.NewDecoder(r.Body)
	response := api.SimpleResponse{Writer: w}
	tokenResponse := api.AdvancedResponse{Writer: w}

	if err := usrDecoder.Decode(&usr); err != nil {
		response.BadRequest("Invalid JSON!")
		return
	}

	// Check password
	if api.InvalidPassword(usr.Password) {
		response.Unauthorized("Invalid password!")
		return
	}

	// Check email
	if api.InvalidEmail(usr.Email) {
		response.BadRequest("Invalid email!")
		return
	}

	// Attempt login
	correct, err := database.Login(usr.Email, usr.Password)

	if errors.Is(err, database.ErrBadPW) {
		response.Unauthorized("Invalid password!")
		return
	} else if errors.Is(err, database.ErrNotFound) {
		response.Unauthorized("Email does not exist!")
		return
	} else if err != nil {
		response.InternalError()
		return
	}

	// Login
	token, err := auth.Create(strconv.Itoa(correct.ID))

	if err != nil {
		response.InternalError()
		return
	}

	tokenResponse.Code = http.StatusOK
	tokenResponse.Data = token
	tokenResponse.Count = 1

	tokenResponse.JSON()
}

func SendPasswordReset(w http.ResponseWriter, r *http.Request) {
	response := api.SimpleResponse{Writer: w}
	emailAddress := r.URL.Query().Get("email")

	// User has verified address
	minfo, err := database.EmailInfo(emailAddress)
	resetToken := strings.Replace(uuid.New().String(), "-", "", -1)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			response.NotFound("Account with email does not exist!")
			return
		}

		response.InternalError()
		return
	}

	if minfo.EmailVerified == 0 {
		response.Unauthorized("This email address is not verified. Password reset unavailable.")
		return
	}

	// Save and send reset token
	if err := cache.SaveResetToken(minfo.Email, resetToken); err != nil {
		response.InternalError()
		return
	}

	m := mail.MailClient{
		To: []string{
			minfo.Email,
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
			"Subject: Password Reset\r\n\r\n"+
			"A password reset has been requested!\n\n%s\r\n",

		config.Data.Mail.User,
		minfo.Email, resetToken,
	))); err != nil {
		log.Error.Println("Could not send email -", err)
		response.InternalError()
		return
	}

	response.Success("Email has been sent!")
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	response := api.SimpleResponse{Writer: w}

	// Validate new password
	password := r.URL.Query().Get("password")

	if api.InvalidPassword(password) {
		response.BadRequest("Invalid password!")
		return
	}

	// Validate reset token
	email, err := cache.GetResetToken(id)

	if err != nil {
		response.InternalError()
		return
	}

	if email == "" {
		response.BadRequest("Invalid token!")
		return
	}

	// Update password and remove token
	if err := database.ResetEmailPassword(email, password); err != nil {
		response.InternalError()
		return
	}

	if err := cache.DeleteResetToken(id); err != nil {
		response.InternalError()
		return
	}

	response.Success("Password has been reset!")
}
