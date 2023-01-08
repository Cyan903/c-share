package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Cyan903/c-share/internal/database"
	"github.com/Cyan903/c-share/pkg/api"
	"github.com/Cyan903/c-share/pkg/log"
)

type account struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var acc account
	accDecoder := json.NewDecoder(r.Body)
	response := api.SimpleResponse{
		Writer: w,
	}

	if err := accDecoder.Decode(&acc); err != nil {
		response.BadRequest(fmt.Sprintf("Could not decode json | %s", err.Error()))
		return
	}

	// Check nickname
	if api.ValidateNickname(acc.Nickname) {
		response.BadRequest("Invalid nickname")
		return
	}

	// Check password
	if api.ValidatePassword(acc.Password) {
		response.BadRequest("Invalid password")
		return
	}

	// Check email
	if api.ValidateEmail(acc.Email) {
		response.BadRequest("Invalid email")
		return
	}

	inUse, err := database.EmailUsed(acc.Email)

	if err != nil {
		response.InternalError()
		return
	}

	if inUse {
		response.Conflict("Email is in use!")
		return
	}

	// Register request
	uid, err := database.Register(acc.Nickname, acc.Email, acc.Password)

	if err != nil {
		response.InternalError()
		return
	}

	log.Info.Println(uid)
}
