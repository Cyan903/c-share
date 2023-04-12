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
		r.Get("/f/{id}", handlers.GetPrivate)
		r.Get("/f/{id}/info", handlers.PrivateFileInfo)
		r.Patch("/f/{id}/edit", handlers.EditFileInfo)

		// TODO: These should be limited
		r.Post("/profile/nickname", handlers.UpdateNickname)
		r.Post("/profile/password", handlers.UpdatePassword)
		r.Post("/profile/email", handlers.UpdateEmail)
		r.Post("/profile/verify", handlers.SendVerification)
		r.Post("/profile/{id}", handlers.VerifyEmail)

		r.Post("/upload", handlers.Upload)
		r.Delete("/upload", handlers.DeleteUpload)
	})

	// mux.Post("/auth/pwreset", handlers.PasswordReset)
	mux.Post("/auth/register", handlers.Register)
	mux.Post("/auth/login", handlers.Login)

	// mux.Get("/f", ?) // server stats (dev only)
	mux.Get("/f/{id}", handlers.GetFile)

	return mux
}
