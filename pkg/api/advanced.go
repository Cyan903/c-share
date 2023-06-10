package api

import (
	"encoding/json"
	"net/http"

	"github.com/Cyan903/c-share/pkg/log"
)

type AdvancedResponse struct {
	Code   int                 `json:"code"`
	Count  int                 `json:"count"`
	Data   any                 `json:"data"`
	Writer http.ResponseWriter `json:"-"`
}

func (s *AdvancedResponse) JSON() {
	js, err := json.Marshal(s)

	if err != nil {
		log.Error.Println("Could not marshal json!")
		return
	}

	s.Writer.Header().Set("Content-Type", "application/json")
	s.Writer.WriteHeader(s.Code)
	s.Writer.Write(js)
}
