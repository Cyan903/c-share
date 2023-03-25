package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Cyan903/c-share/internal/database"
	"github.com/Cyan903/c-share/pkg/api"
	"github.com/Cyan903/c-share/pkg/config"
	"github.com/Cyan903/c-share/pkg/log"
	"github.com/go-chi/chi/v5"
)

func GetFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	pass := r.URL.Query().Get("pass")
	response := api.SimpleResponse{Writer: w}
	file, err := database.GetFile(id, pass)

	if errors.Is(database.ErrNotFound, err) {
		response.NotFound("File not found!")
		return
	} else if errors.Is(database.ErrBadPW, err) {
		response.Unauthorized("Invalid password!")
		return
	} else if err != nil {
		response.InternalError()
		return
	}

	if file.Permissions == 1 {
		response.NotFound("File not found!")
		return
	}

	if file.Permissions == 0 && pass != "" {
		response.BadRequest("Password not required!")
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
