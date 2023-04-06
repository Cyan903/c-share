package handlers

import (
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

// Requires password
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func UpdateEmail(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
