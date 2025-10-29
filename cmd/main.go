package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/luthfiarsyad/mms/config"
	"github.com/luthfiarsyad/mms/internal/interface/http"
)

func main() {
	if err := config.Load(""); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	cfg := config.Get()
	if cfg == nil {
		log.Fatalf("config is nil after load")
	}

	// set gin mode from config
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	http.SetupRoutes(r)

	// run using address from config
	if err := r.Run(cfg.Server.Address); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
