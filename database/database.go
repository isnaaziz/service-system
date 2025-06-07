package database

import (
	"log"
	"os"
	"service_system/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Fatal("JWT_SECRET_KEY is not set in the environment")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	// AutoMigrate
	if err := DB.AutoMigrate(&models.User{}, &models.UserSession{}); err != nil {
		log.Fatal("AutoMigrate failed: ", err)
	}

	// Manual index
	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_user_sessions_userid_isdeleted_createdat ON user_sessions (user_id, is_deleted, created_at)`)
}
