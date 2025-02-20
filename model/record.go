package model

import "time"

type Record struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string    `json:"title" gorm:"not null; default: ''"`
	Artist      string    `json:"artist" gorm:"not null; default: ''"`
	Genre       string    `json:"genre" gorm:"not null; default: ''"`
	Style       string    `json:"style" gorm:"not null; default: ''"`
	ReleaseYear int       `json:"release_year" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	// time.Time 型の場合、nullは扱えない
	// null を許容したい場合は、*time.Time 型を使う
	UpdatedAt *time.Time `json:"updated_at" gorm:"default:null"`
}

type RecordResponse struct {
	// IDはupdateで使うので返す
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Genre       string `json:"genre"`
	Style       string `json:"style"`
	ReleaseYear int    `json:"release_year"`
}
