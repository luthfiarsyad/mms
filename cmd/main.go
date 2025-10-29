package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luthfiarsyad/mms/internal/interface/http"
)

func main() {
	r := gin.New()

	http.SetupRoutes(r)
	
	r.Run(":8080")
}