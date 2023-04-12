package handlers

import (
	"errors"
	"net/http"

	"github.com/Cyan903/c-share/internal/database"
	"github.com/Cyan903/c-share/pkg/api"
	"github.com/golang-jwt/jwt"
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

func SendVerification(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
