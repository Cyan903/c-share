package handlers

import (
	"encoding/json"
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

type fileSearch struct {
	Listing string `json:"listing"`    // [any, public, private, unlisted]
	Type    string `json:"file_type"`  // [any, text/html]
	Order   string `json:"file_order"` // [any, size, type, comment, permission, date]
	Sort    string `json:"sort"`       // [asc, desc]
}

// ? page = 0 & search = ? (optional)
func FilesListing(w http.ResponseWriter, r *http.Request) {
	var fsearch fileSearch

	id := r.Context().Value(jwt.StandardClaims{}).(*jwt.StandardClaims)
	response := api.SimpleResponse{Writer: w}
	jsonResponse := api.AdvancedResponse{Writer: w}
	fdecoder := json.NewDecoder(r.Body)

	if err := fdecoder.Decode(&fsearch); err != nil {
		response.BadRequest("Invalid JSON!")
		return
	}

	search := r.URL.Query().Get("search")
	page := r.URL.Query().Get("page")
	pagn, err := strconv.Atoi(page)

	if err != nil || pagn < 0 {
		response.BadRequest("Invalid page!")
		return
	}

	// Validate params
	if !slices.Contains([]string{"any", "public", "private", "unlisted"}, fsearch.Listing) {
		response.BadRequest("Invalid listing!")
		return
	}

	if fsearch.Type == "" {
		response.BadRequest("Invalid file type!")
		return
	}

	if !slices.Contains([]string{"any", "size", "type", "comment", "permission", "date"}, fsearch.Order) {
		response.BadRequest("Invalid order!")
		return
	}

	if fsearch.Sort != "asc" && fsearch.Sort != "desc" {
		response.BadRequest("Invalid sort!")
		return
	}

	if len(search) > 99 {
		response.BadRequest("Invalid search!")
		return
	}

	files, count, err := database.FileListing(id.Issuer, pagn, fsearch.Listing, fsearch.Type, fsearch.Order, fsearch.Sort, search)

	if err != nil {
		response.InternalError()
		return
	}

	jsonResponse.Code = http.StatusOK
	jsonResponse.Data = files
	jsonResponse.Count = count

	jsonResponse.JSON()
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
