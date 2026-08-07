package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       UUIDBinary `json:"id" gorm:"type:binary(16);primaryKey;not null"`
	Username string     `json:"username" gorm:"unique;not null"`
	Name     string     `json:"name" gorm:"not null"`
	Email    string     `json:"email" gorm:"unique;not null"`
	Password string     `json:"password" gorm:"not null"`
	Roles    []Role     `json:"roles" gorm:"many2many:user_roles;"`
}

// Hook GORM: Otomatis generate UUID v4 baru jika ID masih kosong
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if uuid.UUID(u.ID) == uuid.Nil {
		u.ID = UUIDBinary(uuid.New())
	}
	return
}