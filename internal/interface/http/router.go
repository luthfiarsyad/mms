package http

import (
	"github.com/gin-gonic/gin"
	"github.com/luthfiarsyad/mms/internal/interface/http/handler"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	v1 := api.Group("/v1")

	// --- AUTH ROUTES ---
	authHandler := handler.NewAuthHandler()
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// --- USERS ROUTES ---
	users := v1.Group("/users")
	{
		users.POST("", createUser)
		users.GET("", listUsers)
		users.GET("/:id", getUser)
		users.PUT("/:id", updateUser)
		users.DELETE("/:id", deleteUser)
	}

	// --- TRANSACTIONS ROUTES ---
	tx := v1.Group("/transactions")
	{
		tx.POST("", createTransaction)
		tx.GET("", listTransactions)
		tx.GET("/:id", getTransaction)
	}
}

func createUser(c *gin.Context)        { c.JSON(501, gin.H{"error": "not implemented"}) }
func listUsers(c *gin.Context)         { c.JSON(501, gin.H{"error": "not implemented"}) }
func getUser(c *gin.Context)           { c.JSON(501, gin.H{"error": "not implemented"}) }
func updateUser(c *gin.Context)        { c.JSON(501, gin.H{"error": "not implemented"}) }
func deleteUser(c *gin.Context)        { c.JSON(501, gin.H{"error": "not implemented"}) }
func createTransaction(c *gin.Context) { c.JSON(501, gin.H{"error": "not implemented"}) }
func listTransactions(c *gin.Context)  { c.JSON(501, gin.H{"error": "not implemented"}) }
func getTransaction(c *gin.Context)    { c.JSON(501, gin.H{"error": "not implemented"}) }
