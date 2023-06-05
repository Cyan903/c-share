package handlers

import (
	"net/http"

	"github.com/Cyan903/c-share/pkg/api"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	response := api.SimpleResponse{Writer: w}

	response.NotFound("Route not found!")
}

func RateLimit(w http.ResponseWriter, r *http.Request) {
	response := api.SimpleResponse{Writer: w}
	
	response.Code = http.StatusTooManyRequests
	response.Message = "Too many requests!"
	
	response.JSON()
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	response := api.SimpleResponse{Writer: w}

	response.Code = http.StatusMethodNotAllowed
	response.Message = "Method not allowed!"

	response.JSON()
}
