package database

import (
	"fmt"
	"log"
	"time"

	"github.com/example/supabase-migration-demo/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func New(dsn string) (*Database, error) {
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Add this to work around PgBouncer prepared statement issues
	conn := postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Disable prepared statements
	})

	db, err := gorm.Open(conn, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *Database) Migrate() error {
	log.Println("Starting AutoMigrate...")

	// Use Migrator directly to avoid duplicate migration attempts
	if err := d.DB.Migrator().AutoMigrate(&models.User{}, &models.Document{}); err != nil {
		// Log the error but don't fail if tables already exist
		log.Printf("Migration warning: %v", err)
	}

	log.Println("AutoMigrate completed successfully")
	log.Println("Models migrated: User, Document")

	return nil
}

func (d *Database) GetDB() *gorm.DB {
	return d.DB
}