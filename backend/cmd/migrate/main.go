package main

import (
	"log"
	"os"
	"todo-api/infra/gorm/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Todo{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("migration completed")
}
