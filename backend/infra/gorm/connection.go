package gorm

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			TranslateError: true, //ErrDuplicatedKeyを使うために必要　参考：https://gorm.io/docs/error_handling.html#Dialect-Translated-Errors
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	return db, nil
}
