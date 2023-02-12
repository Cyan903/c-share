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
	})

	mux.Post("/auth/register", handlers.Register)
	mux.Post("/auth/login", handlers.Login)

	return mux
}
