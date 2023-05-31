package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Cyan903/c-share/internal/database"
	"github.com/Cyan903/c-share/pkg/api"
	"github.com/Cyan903/c-share/pkg/config"
	"github.com/Cyan903/c-share/pkg/log"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"golang.org/x/exp/slices"
)

// ? page = 0
// & listing = [any, public, private, unlisted]
// & type = [any, text/html]
// & order = [any, size, type, comment, permission, date]
// & sort = [asc, desc]
// & search = ? (optional)
func FilesListing(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)
	response := api.SimpleResponse{Writer: w}
	json := api.AdvancedResponse{Writer: w}

	page := r.URL.Query().Get("page")
	pagn, err := strconv.Atoi(page)

	if err != nil || pagn < 0 {
		response.BadRequest("Invalid page!")
		return
	}

	perms := []string{"any", "public", "private", "unlisted"}
	perm := r.URL.Query().Get("listing")

	if !slices.Contains(perms, perm) {
		response.BadRequest("Invalid listing!")
		return
	}

	fileType := r.URL.Query().Get("type")

	if fileType == "" {
		response.BadRequest("Invalid file type!")
		return
	}

	orders := []string{"size", "type", "comment", "permission", "date"}
	order := r.URL.Query().Get("order")

	if !slices.Contains(orders, order) {
		response.BadRequest("Invalid order!")
		return
	}

	sort := r.URL.Query().Get("sort")

	if sort != "asc" && sort != "desc" {
		response.BadRequest("Invalid sort!")
		return
	}

	search := r.URL.Query().Get("search")

	if len(search) > 99 {
		response.BadRequest("Invalid search!")
		return
	}

	files, count, err := database.FileListing(id.Issuer, pagn, perm, fileType, order, sort, search)

	if err != nil {
		response.InternalError()
		return
	}

	json.Code = http.StatusOK
	json.Data = files
	json.Count = count

	json.JSON()
}

func GetPrivate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	response := api.SimpleResponse{Writer: w}
	uid := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)
	file, err := database.GetPrivateFile(id, uid.Issuer)

	if errors.Is(database.ErrNotFound, err) {
		response.NotFound("File not found!")
		return
	} else if err != nil {
		response.InternalError()
		return
	}

	p := fmt.Sprintf("%s/%s", config.Data.UploadPath, file.ID)

	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		log.Warning.Println("File exists in DB, but not on disk!", p)
		response.NotFound("File not found!")

		return
	}

	http.ServeFile(w, r, p)
}

func PrivateFileInfo(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)
	file := chi.URLParam(r, "id")

	response := api.SimpleResponse{Writer: w}
	jsonResponse := api.AdvancedResponse{Writer: w}

	info, err := database.FileInfo(id.Issuer, file)

	if errors.Is(database.ErrNotFound, err) {
		response.NotFound("File not found!")
		return
	} else if err != nil {
		response.InternalError()
		return
	}

	jsonResponse.Code = 200
	jsonResponse.Count = 6
	jsonResponse.Data = info

	jsonResponse.JSON()
}

func EditFileInfo(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)
	file := chi.URLParam(r, "id")
	response := api.SimpleResponse{Writer: w}

	password := r.URL.Query().Get("password")
	comment := r.URL.Query().Get("comment")
	upriv, priv := r.URL.Query().Get("perm"), 0

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

	if api.InvalidFilename(comment) {
		response.BadRequest("Invalid file comment!")
		return
	}

	if priv == 2 && api.InvalidPassword(password) {
		response.BadRequest("Invalid password!")
		return
	}

	if upriv != "unlisted" && password != "" {
		response.BadRequest("Cannot have password on public/private files!")
		return
	} else if upriv == "unlisted" && password == "" {
		response.BadRequest("Password required for unlisted files!")
		return
	}

	err := database.UpdateFileInfo(file, id.Issuer, password, comment, priv)

	if errors.Is(database.ErrNotFound, err) {
		response.NotFound("File not found!")
		return
	} else if err != nil {
		response.InternalError()
		return
	}

	response.Success("File has been updated!")
}
