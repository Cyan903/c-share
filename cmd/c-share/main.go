package main

import (
	"os"

	"github.com/Cyan903/c-share/internal/database"
	"github.com/Cyan903/c-share/internal/router"
	"github.com/Cyan903/c-share/pkg/config"
)

func main() {
	cfg := config.LoadConfig()
	conn, err := database.OpenDB(cfg.DSN)

	if err != nil {
		os.Exit(1)
	}

	database.Conn = conn
	config.Dev = cfg.Mode == "development"

	router.Serve(cfg.Port)
}
