package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	response := api.SimpleResponse{Writer: w}
	successResponse := api.AdvancedResponse{Writer: w}

	// Check file size
	if r.ContentLength > config.Data.UploadLimit {
		response.BadRequest("Too large!")
		return
	}

	r.ParseMultipartForm(config.Data.UploadLimit)
	r.Body = http.MaxBytesReader(w, r.Body, config.Data.UploadLimit)

	// Accept request
	upriv, priv := r.URL.Query().Get("perm"), 0
	pass := r.URL.Query().Get("password")
	comment := r.URL.Query().Get("comment")

	id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)
	file, handler, err := r.FormFile("upload")

	if api.InvalidFilename(comment) {
		response.BadRequest("Invalid file comment!")
		return
	}

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

	// Confirm password
	if pass != "" && priv != 2 {
		response.BadRequest("Password not required!")
		return
	}

	if priv == 2 && api.InvalidPassword(pass) {
		response.BadRequest("Invalid password!")
		return
	}

	// Create ID
	rid, err := database.RandomID()

	if err != nil {
		response.InternalError()
		return
	}

	// Upload file
	tfile, err := os.Create(config.Data.UploadPath + "/" + rid)

	if err != nil {
		response.InternalError()
		log.Error.Println(err)
		return
	}

	defer tfile.Close()
	fbytes, err := io.ReadAll(file)

	if err != nil {
		response.InternalError()
		log.Error.Println(err)
		return
	}

	tfile.Write(fbytes)

	// Upload to database
	if err := database.UploadFile(rid, id.Issuer, handler.Size, handler.Header.Get("Content-Type"), pass, comment, priv); err != nil {
		response.InternalError()
		log.Error.Println("Could not add file to database -", err)
		return
	}

	// Update Storage
	storage, err := database.UpdateStorage(id.Issuer)

	if err != nil {
		response.InternalError()
		return
	}

	successResponse.Code = 200
	successResponse.Count = 2
	successResponse.Data = struct {
		ID      string `json:"id"`
		Storage string `json:"storage"`
	}{rid, strconv.FormatInt(storage, 10)}

	successResponse.JSON()
}

func DeleteUpload(w http.ResponseWriter, r *http.Request) {
	var files []string

	id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)
	fileDecoder := json.NewDecoder(r.Body)
	response := api.SimpleResponse{Writer: w}

	if err := fileDecoder.Decode(&files); err != nil {
		response.BadRequest("Invalid JSON!")
		return
	}

	if len(files) <= 0 {
		response.BadRequest("No files?")
		return
	}

	// Check dups
	sarg := make(map[string]bool)

	for _, v := range files {
		if _, has := sarg[v]; has {
			response.BadRequest("Duplicate IDs!")
			return
		}

		sarg[v] = false
	}

	// Does user own files / Do they exist?
	uid, err := strconv.Atoi(id.Issuer)

	if err != nil {
		response.InternalError()
		log.Error.Println("Could not convert user ID -", err)
		return
	}

	o, err := database.OwnFiles(files, uid)

	if err != nil {
		response.InternalError()
		log.Error.Println("Could not delete files -", err)
		return
	}

	if len(o) != 0 {
		response.Conflict("Invalid IDs: " + strings.Join(o, ", "))
		return
	}

	// Remove from database/disk
	if err := database.DeleteFiles(id.Issuer, files); err != nil {
		response.InternalError()
		log.Error.Println("Could not remove file from DB -", err)
		return
	}

	for _, f := range files {
		if err := os.Remove(config.Data.UploadPath + "/" + f); err != nil {
			response.InternalError()
			log.Error.Println("Could not remove file from disk -", err)
			return
		}
	}

	// Update Storage
	storage, err := database.UpdateStorage(id.Issuer)

	if err != nil {
		response.InternalError()
		return
	}

	response.Success(strconv.FormatInt(storage, 10))
}
