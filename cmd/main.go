package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/supabase-migration-demo/internal/config"
	"github.com/example/supabase-migration-demo/internal/database"
	"github.com/example/supabase-migration-demo/internal/storage"
	"github.com/example/supabase-migration-demo/repositories"
	"github.com/example/supabase-migration-demo/services"
	"github.com/example/supabase-migration-demo/controllers"
	"github.com/example/supabase-migration-demo/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting Go + PostgreSQL + S3 Demo Application")

	cfg := config.Load()

	db, err := database.New(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	storageClient, err := storage.New(storage.Config{
		ProjectID: cfg.S3ProjectID,
		Endpoint:  cfg.S3Endpoint,
		Region:    cfg.S3Region,
		AccessKey: cfg.S3AccessKey,
		SecretKey: cfg.S3SecretKey,
		Bucket:    cfg.S3Bucket,
	})
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	userRepo := repositories.NewUserRepository(db.GetDB())
	docRepo := repositories.NewDocumentRepository(db.GetDB())

	userService := services.NewUserService(userRepo)
	docService := services.NewDocumentService(docRepo, storageClient)

	userCtrl := controllers.NewUserController(userService)
	docCtrl := controllers.NewDocumentController(docService, storageClient)

	router := gin.Default()
	router.MaxMultipartMemory = 10 << 20

	routes.Setup(router, userCtrl, docCtrl)

	go func() {
		log.Printf("Server starting on %s", cfg.GetServerAddress())
		if err := router.Run(cfg.GetServerAddress()); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}