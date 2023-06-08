package handlers

import (
	"net/http"

	"github.com/Cyan903/c-share/internal/database"
	"github.com/Cyan903/c-share/pkg/api"
	"github.com/golang-jwt/jwt"
)

func APIEmailCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)
		response := api.SimpleResponse{Writer: w}
		abt, err := database.About(id.Issuer)

		if err != nil {
			response.InternalError()
			return
		}

		if abt.EmailVerified == 0 {
			response.Unauthorized("You must have a verified email to interact with the API!")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetAPIToken(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)
	response := api.SimpleResponse{Writer: w}
	jsonRes := api.AdvancedResponse{Writer: w}

	token, err := database.UserAPIToken(id.Issuer)

	if err != nil {
		response.InternalError()
		return
	}

	jsonRes.Code = http.StatusOK
	jsonRes.Count = 2
	jsonRes.Data = token

	jsonRes.JSON()
}

func GenerateAPIToken(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)
	response := api.SimpleResponse{Writer: w}

	abt, err := database.About(id.Issuer)
	utoken, uerr := database.UserAPIToken(id.Issuer)

	if err != nil || uerr != nil {
		response.InternalError()
		return
	}

	// Remove token if it exists
	if utoken.Token != "" {
		if err := database.DeleteAPIToken(id.Issuer); err != nil {
			response.InternalError()
			return
		}
	}

	token, err := database.GenerateToken(abt.Email)

	if err != nil {
		response.InternalError()
		return
	}

	if err := database.AddAPIToken(id.Issuer, token); err != nil {
		response.InternalError()
		return
	}

	response.Success(token)
}

func RemoveAPIToken(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)
	response := api.SimpleResponse{Writer: w}

	token, err := database.UserAPIToken(id.Issuer)

	if err != nil {
		response.InternalError()
		return
	}

	if token.Token == "" {
		response.BadRequest("You don't have an API token!")
		return
	}

	if err := database.DeleteAPIToken(id.Issuer); err != nil {
		response.InternalError()
		return
	}

	response.Success("Token as been removed!")
}
