package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/luthfiarsyad/mms/config"
	"github.com/luthfiarsyad/mms/internal/infrastructure/logger"
	"github.com/luthfiarsyad/mms/internal/infrastructure/persistence/mysql"
	"github.com/luthfiarsyad/mms/internal/infrastructure/http/middleware"
	"github.com/luthfiarsyad/mms/internal/interface/http"
)

func main() {
	// --- Load config ---
	if err := config.Load(""); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	cfg := config.Get()
	if cfg == nil {
		log.Fatalf("config is nil after load")
	}

	// --- Initialize Zerolog ---
	logger.Init(cfg.Log.Level)
	log := logger.Get()
	log.Info().Msg("Logger initialized")

	// --- Initialize Database ---
	if _, err := mysql.Connect(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	log.Info().Msg("Database connected")

	// --- Setup Gin ---
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	// Apply middlewares
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())

	// --- Setup routes ---
	http.SetupRoutes(r)

	// --- Run server ---
	addr := cfg.Server.Address
	if addr == "" {
		addr = ":8080"
	}
	log.Info().Msgf("Server starting on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
