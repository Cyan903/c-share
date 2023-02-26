package router

import (
	"net/http"

	"github.com/Cyan903/c-share/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func routes() http.Handler {
	// TODO: CORS
	// TODO: Rate limit
	mux := chi.NewRouter()

	mux.Route("/@me", func(r chi.Router) {
		r.Use(handlers.TokenCheck)
		r.Get("/", handlers.WhoAmI)

		r.Get("/f", handlers.FilesListing)
		// r.Get("/f/{id}", ?) // private files
		// r.Get("/f/{id}/info", ?) // file info/type

		r.Post("/upload", handlers.Upload)
		r.Delete("/upload", handlers.DeleteUpload)
	})

	mux.Post("/auth/register", handlers.Register)
	mux.Post("/auth/login", handlers.Login)
	
	// mux.Get("/f", ?) // server stats (dev only)
	mux.Get("/f/{id}", handlers.GetFile)

	return mux
}
