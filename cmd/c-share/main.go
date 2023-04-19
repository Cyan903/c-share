package main

import (
	"os"

	"github.com/Cyan903/c-share/internal/cache"
	"github.com/Cyan903/c-share/internal/database"
	"github.com/Cyan903/c-share/internal/router"
	"github.com/Cyan903/c-share/pkg/config"
)

func main() {
	cfg := config.LoadConfig()
	conn, err := database.OpenDB(cfg.DSN)
	cerr := cache.OpenCache(cfg.Cache.Address, cfg.Cache.Password, cfg.Cache.DB)

	if err != nil || cerr != nil {
		os.Exit(1)
	}

	database.Conn = conn
	router.Serve(cfg.Port)
}
