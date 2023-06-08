package handlers

import (
	"context"
	"net/http"

	"github.com/Cyan903/c-share/internal/database"
	"github.com/Cyan903/c-share/pkg/api"
	"github.com/Cyan903/c-share/pkg/auth"
	"github.com/golang-jwt/jwt"
)

func TokenCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Token")
		response := api.SimpleResponse{Writer: w}

		if len(token) == 0 {
			response.Unauthorized("No token provided!")
			return
		}

		check, err := auth.VerifyToken(token)

		if err != nil {
			response.Unauthorized("Invalid token!")
			return
		}

		claims := check.Claims
		ctx := context.WithValue(r.Context(), jwt.StandardClaims{}, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WhoAmI(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)
	response := api.SimpleResponse{Writer: w}
	json := api.AdvancedResponse{Writer: w}
	abt, err := database.About(id.Issuer)

	if err != nil {
		response.InternalError()
		return
	}

	json.Code = http.StatusOK
	json.Count = 5
	json.Data = abt

	json.JSON()
}

func Upload(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)

	FileUpload(id.Issuer, w, r)
}

func DeleteUpload(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)

	FileDelete(id.Issuer, w, r)
}
