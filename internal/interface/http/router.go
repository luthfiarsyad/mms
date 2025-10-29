package http

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine) {
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	api := r.Group("/api")
	v1 := api.Group("/v1")

	users := v1.Group("/users")
	{
		users.POST("", createUser)
		users.GET("", listUsers)
		users.GET("/:id", getUser)
		users.PUT("/:id", updateUser)
		users.DELETE("/:id", deleteUser)
	}

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
