package entity

import "time"

type Base struct {
	ID        string     `gorm:"id" json:"id"`
	CreatedAt *time.Time `gorm:"created_at" json:"created-at"`
	UpdatedAt *time.Time `gorm:"updated_at" json:"updated-at"`
	DeletedAt *time.Time `gorm:"deleted_at" json:"deleted-at"`
}
