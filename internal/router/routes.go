package router

import (
	"net/http"
	"time"

	"github.com/Cyan903/c-share/internal/handlers"
	"github.com/Cyan903/c-share/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	info := httprate.Limit(10, 10*time.Minute, httprate.WithLimitHandler(handlers.RateLimit))
	auth := httprate.Limit(15, 10*time.Minute, httprate.WithLimitHandler(handlers.RateLimit))
	email := httprate.Limit(3, 5*time.Minute, httprate.WithLimitHandler(handlers.RateLimit))
	reset := httprate.Limit(10, 5*time.Minute, httprate.WithLimitHandler(handlers.RateLimit))

	mux.Use(httprate.Limit(250, 1*time.Minute, httprate.WithLimitHandler(handlers.RateLimit)))
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   config.Data.CorsAllow,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Route("/@me", func(r chi.Router) {
		r.Use(handlers.TokenCheck)
		r.Get("/", handlers.WhoAmI)

		r.Get("/f", handlers.FilesListing)
		r.Get("/f/{id}", handlers.GetPrivate)
		r.Get("/f/{id}/info", handlers.PrivateFileInfo)
		r.Post("/f/{id}/edit", handlers.EditFileInfo)

		r.With(info).Post("/profile/nickname", handlers.UpdateNickname)
		r.With(info).Post("/profile/password", handlers.UpdatePassword)
		r.With(info).Post("/profile/email", handlers.UpdateEmail)

		r.With(email).Post("/profile/verify", handlers.SendVerification)
		r.With(email).Post("/profile/{id}", handlers.VerifyEmail)

		r.Post("/upload", handlers.Upload)
		r.Delete("/upload", handlers.DeleteUpload)
	})

	mux.With(auth).Post("/auth/register", handlers.Register)
	mux.With(auth).Post("/auth/login", handlers.Login)

	mux.With(reset).Post("/auth/pwreset", handlers.SendPasswordReset)
	mux.With(reset).Post("/auth/{id}", handlers.ResetPassword)

	mux.Get("/f", handlers.ServerFiles)
	mux.Get("/f/{id}", handlers.GetFile)

	mux.Get("/", handlers.ServerStats)
	mux.NotFound(handlers.NotFound)
	mux.MethodNotAllowed(handlers.MethodNotAllowed)

	return mux
}
