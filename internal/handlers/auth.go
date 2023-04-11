package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Cyan903/c-share/internal/database"
	"github.com/Cyan903/c-share/pkg/api"
	"github.com/Cyan903/c-share/pkg/auth"
)

type account struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type user struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var acc account

	accDecoder := json.NewDecoder(r.Body)
	response := api.SimpleResponse{Writer: w}
	tokenResponse := api.AdvancedResponse{Writer: w}

	if err := accDecoder.Decode(&acc); err != nil {
		response.BadRequest("Invalid JSON!")
		return
	}

	// Check nickname
	if api.InvalidNickname(acc.Nickname) {
		response.BadRequest("Invalid nickname!")
		return
	}

	// Check password
	if api.InvalidPassword(acc.Password) {
		response.BadRequest("Invalid password!")
		return
	}

	// Check email
	if api.InvalidEmail(acc.Email) {
		response.BadRequest("Invalid email!")
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

	// Login
	token, err := auth.Create(strconv.Itoa(int(uid)))

	if err != nil {
		response.InternalError()
		return
	}

	tokenResponse.Code = http.StatusOK
	tokenResponse.Data = token
	tokenResponse.Count = 1

	tokenResponse.JSON()
}

func Login(w http.ResponseWriter, r *http.Request) {
	var usr user

	usrDecoder := json.NewDecoder(r.Body)
	response := api.SimpleResponse{Writer: w}
	tokenResponse := api.AdvancedResponse{Writer: w}

	if err := usrDecoder.Decode(&usr); err != nil {
		response.BadRequest("Invalid JSON!")
		return
	}

	// Check password
	if api.InvalidPassword(usr.Password) {
		response.Unauthorized("Invalid password!")
		return
	}

	// Check email
	if api.InvalidEmail(usr.Email) {
		response.BadRequest("Invalid email!")
		return
	}

	// Attempt login
	correct, err := database.Login(usr.Email, usr.Password)

	if errors.Is(err, database.ErrBadPW) {
		response.Unauthorized("Invalid password!")
		return
	} else if errors.Is(err, database.ErrNotFound) {
		response.Unauthorized("Email does not exist!")
		return
	} else if err != nil {
		response.InternalError()
		return
	}

	// Login
	token, err := auth.Create(strconv.Itoa(correct.ID))

	if err != nil {
		response.InternalError()
		return
	}

	tokenResponse.Code = http.StatusOK
	tokenResponse.Data = token
	tokenResponse.Count = 1

	tokenResponse.JSON()
}
