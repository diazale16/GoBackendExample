package routes

import (
	"github.com/example/supabase-migration-demo/controllers"
	"github.com/gin-gonic/gin"
)

func Setup(
	router *gin.Engine,
	userCtrl *controllers.UserController,
	docCtrl *controllers.DocumentController,
) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	users := router.Group("/users")
	{
		users.POST("", userCtrl.Create)
		users.GET("", userCtrl.GetAll)
		users.GET("/:id", userCtrl.GetByID)
		users.DELETE("/:id", userCtrl.Delete)
	}

	documents := router.Group("/documents")
	{
		documents.POST("/upload", docCtrl.Upload)
		documents.GET("", docCtrl.GetAll)
		documents.GET("/:id", docCtrl.GetByID)
		documents.DELETE("/:id", docCtrl.Delete)
	}
}