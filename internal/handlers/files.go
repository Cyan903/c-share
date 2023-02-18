package handlers

import (
	"database/sql"
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
	response := api.SimpleResponse{Writer: w}
	file, err := database.GetFile(id)

	if errors.Is(sql.ErrNoRows, err) {
		response.NotFound("File not found!")
		return
	} else if err != nil {
		response.InternalError()
		return
	}

	if file.Permissions == 1 {
		response.NotFound("File not found!")
		return
	} else if file.Permissions == 2 {
		// TODO
		log.Info.Println("Unlisted file hit")
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
