package main

import (
	"fmt"
	"log"
	"record-shop-rest-api/db"
	"record-shop-rest-api/model"
)

func main() {
	dbConn := db.NewDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	// DBに反映させたいModel構造を渡す
	// {}でフィールドの値を0値にしている
	err := dbConn.AutoMigrate(&model.User{}, &model.Record{}, &model.Detail{}, &model.Track{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	err = dbConn.Exec(`
		ALTER TABLE records
		ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;
	`).Error
	if err != nil {
		log.Fatalf("failed to add default timestamps: %v", err)
	}
}
