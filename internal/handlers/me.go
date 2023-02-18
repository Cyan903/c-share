package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Cyan903/c-share/internal/database"
	"github.com/Cyan903/c-share/pkg/api"
	"github.com/Cyan903/c-share/pkg/auth"
	"github.com/Cyan903/c-share/pkg/config"
	"github.com/Cyan903/c-share/pkg/log"
	"github.com/golang-jwt/jwt"
)

func TokenCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Token")
		response := api.SimpleResponse{Writer: w}

		if len(token) == 0 {
			response.Unauthorized("No token provided")
			return
		}

		check, err := auth.VerifyToken(token)

		if err != nil {
			response.Unauthorized("Invalid token")
			return
		}

		claims := check.Claims
		ctx := context.WithValue(r.Context(), jwt.StandardClaims{}, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// TODO: Improve
func WhoAmI(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)
	response := api.SimpleResponse{Writer: w}
	abt, err := database.About(id.Issuer)

	if err != nil {
		response.InternalError()
		return
	}

	response.Success(fmt.Sprintf("%d, %s, %s", abt.ID, abt.Nickname, abt.Register))
}

func Upload(w http.ResponseWriter, r *http.Request) {
	response := api.SimpleResponse{Writer: w}

	// Check file size
	if r.ContentLength > config.Data.UploadLimit {
		response.BadRequest("Too large!")
		return
	}

	r.ParseMultipartForm(config.Data.UploadLimit)
	r.Body = http.MaxBytesReader(w, r.Body, config.Data.UploadLimit)

	// Accept request
	upriv, priv := r.URL.Query().Get("perm"), 0
	id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)
	file, handler, err := r.FormFile("upload")

	// Accept request
	switch upriv {
	case "public":
		priv = 0
	case "private":
		priv = 1
	case "unlisted":
		priv = 2
	default:
		response.BadRequest("Invalid permission!")
		return
	}

	if err != nil {
		response.BadRequest("File missing!")
		return
	}

	defer file.Close()

	if handler.Size > config.Data.UploadLimit {
		response.BadRequest("File too large!")
		return
	}

	// Create ID
	rid, err := database.RandomID()

	if err != nil {
		response.InternalError()
		return
	}

	// Upload file
	tfile, err := ioutil.TempFile(config.Data.UploadPath, rid)

	if err != nil {
		response.InternalError()
		log.Error.Println(err)
		return
	}

	defer tfile.Close()
	fbytes, err := ioutil.ReadAll(file)

	if err != nil {
		response.InternalError()
		log.Error.Println(err)
		return
	}

	tfile.Write(fbytes)

	// Upload to database
	if err := database.UploadFile(rid, id.Issuer, handler.Size, handler.Header.Get("Content-Type"), priv); err != nil {
		response.InternalError()
		log.Error.Println("Could not add file to database", err)
		return
	}

	response.Success("Uploaded")
}
