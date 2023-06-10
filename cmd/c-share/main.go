package main

import (
	"os"

	"github.com/Cyan903/c-share/internal/cache"
	"github.com/Cyan903/c-share/internal/database"
	"github.com/Cyan903/c-share/internal/router"
	"github.com/Cyan903/c-share/pkg/config"
	"github.com/Cyan903/c-share/pkg/log"
)

func main() {
	cfg := config.LoadConfig()
	conn, err := database.OpenDB(cfg.DSN)
	cerr := cache.OpenCache(cfg.Cache.Address, cfg.Cache.Password, cfg.Cache.DB)

	if err != nil || cerr != nil {
		log.Error.Println(err, cerr)
		os.Exit(1)
	}

	database.Conn = conn

	if err := router.Serve(cfg.Port); err != nil {
		log.Error.Println(err)
		os.Exit(1)
	}
}
