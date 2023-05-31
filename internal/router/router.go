package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Cyan903/c-share/pkg/api"
	"github.com/Cyan903/c-share/pkg/log"
)

func Serve(port int) error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	log.Info.Printf("Running on :%d", port)
	return srv.ListenAndServe()
}

func RateLimit(w http.ResponseWriter, r *http.Request) {
	(&api.SimpleResponse{Writer: w}).TooManyRequests()
}
