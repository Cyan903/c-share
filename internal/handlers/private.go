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

func FilesListing(w http.ResponseWriter, r *http.Request) {
	// page = 0
	// listing = any, public, private, unlisted
	// type = any, text/html
	// order = any, size, type, permission, date
	// sort = asc/desc

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

	orders := []string{"any", "size", "type", "permission", "date"}
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

	files, count, err := database.FileListing(id.Issuer, pagn, perm, fileType, order, sort)

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
		log.Info.Println(err)
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
