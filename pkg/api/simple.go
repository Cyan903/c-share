package api

import (
	"encoding/json"
	"net/http"

	"github.com/Cyan903/c-share/pkg/config"
	"github.com/Cyan903/c-share/pkg/log"
)

type SimpleResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Writer  http.ResponseWriter `json:"-"`
}

func (s *SimpleResponse) JSON() {
	js, err := json.Marshal(s)

	if err != nil {
		log.Error.Println("Could not marshal json!")
		return
	}

	s.Writer.Header().Set("Content-Type", "application/json")
	s.Writer.WriteHeader(s.Code)
	s.Writer.Write(js)
}

// Dev only message
func (s *SimpleResponse) BadRequest(msg string) {
	s.Code = http.StatusBadRequest
	s.Message = "Bad request!"

	if config.Dev {
		s.Message = msg
	}

	s.JSON()
}

func (s *SimpleResponse) Unauthorized(msg string) {
	s.Code = http.StatusUnauthorized
	s.Message = "Unauthorized!"

	if config.Dev {
		s.Message = msg
	}

	s.JSON()
}

// No message
func (s *SimpleResponse) InternalError() {
	s.Code = http.StatusInternalServerError
	s.Message = "Internal server error!"

	s.JSON()
}

// Message required!
func (s *SimpleResponse) NotFound(msg string) {
	s.Code = http.StatusNotFound
	s.Message = msg

	s.JSON()
}

func (s *SimpleResponse) Conflict(msg string) {
	s.Code = http.StatusConflict
	s.Message = msg

	s.JSON()
}

func (s *SimpleResponse) Success(msg string) {
	s.Code = http.StatusOK
	s.Message = msg

	s.JSON()
}

