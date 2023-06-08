package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Cyan903/c-share/internal/cache"
	"github.com/Cyan903/c-share/internal/database"
	"github.com/Cyan903/c-share/pkg/api"
	"github.com/Cyan903/c-share/pkg/config"
	"github.com/Cyan903/c-share/pkg/log"
	"github.com/Cyan903/c-share/pkg/web"
	"github.com/go-chi/chi/v5"
)

func GetFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	pass := r.URL.Query().Get("password")
	jsOnly := r.URL.Query().Get("json") == "true"

	response := api.SimpleResponse{Writer: w}
	file, err := database.GetFile(id, pass)

	if errors.Is(database.ErrNotFound, err) {
		if jsOnly {
			response.NotFound("File not found!")
			return
		}

		w.WriteHeader(http.StatusNotFound)
		web.ParseTemplate(w, "404", id)
		return
	} else if errors.Is(database.ErrBadPW, err) {
		if jsOnly {
			response.Unauthorized("Invalid password!")
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		web.ParseTemplate(w, "password", id)
		return
	} else if err != nil {
		response.InternalError()
		return
	}

	if file.Permissions == 1 {
		if jsOnly {
			response.NotFound("File not found!")
			return
		}

		w.WriteHeader(http.StatusNotFound)
		web.ParseTemplate(w, "404", id)
		return
	}

	p := fmt.Sprintf("%s/%s", config.Data.UploadPath, file.ID)

	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		log.Warning.Println("File exists in DB, but not on disk!", p)

		if jsOnly {
			response.NotFound("File not found!")
			return
		}

		w.WriteHeader(http.StatusNotFound)
		web.ParseTemplate(w, "404", id)
		return
	}

	http.ServeFile(w, r, p)
}

// Dev mode only
func ServerFiles(w http.ResponseWriter, r *http.Request) {
	response := api.SimpleResponse{Writer: w}
	jsonResponse := api.AdvancedResponse{Writer: w}

	if config.Data.Mode != "development" {
		response.NotFound("Route not found!")
		return
	}

	info, err := database.ServerStorageInfo()

	if err != nil {
		response.InternalError()
		return
	}

	jsonResponse.Code = 200
	jsonResponse.Count = 3
	jsonResponse.Data = info

	jsonResponse.JSON()
}

// Ok to expose to public
func ServerStats(w http.ResponseWriter, r *http.Request) {
	response := api.SimpleResponse{Writer: w}
	jsonResponse := api.AdvancedResponse{Writer: w}

	statsCache, success, err := cache.GetServerInfo()

	if err != nil {
		response.InternalError()
		return
	}

	if success {
		jsonResponse.Code = 200
		jsonResponse.Count = 3
		jsonResponse.Data = statsCache

		jsonResponse.JSON()
		return
	}

	stats, err := database.ServerStatsInfo()

	if err != nil {
		response.InternalError()
		return
	}

	if err := cache.SaveServerStats(stats.Users, stats.Storage, stats.Total); err != nil {
		response.InternalError()
		return
	}

	jsonResponse.Code = 200
	jsonResponse.Count = 3
	jsonResponse.Data = stats

	jsonResponse.JSON()
}
