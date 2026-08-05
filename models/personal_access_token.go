package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PersonalAccessToken struct {
	ID        UUIDBinary `json:"id" gorm:"type:binary(16);primaryKey"`
	UserID    UUIDBinary `json:"user_id" gorm:"type:binary(16);not null"`
	User      User       `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TokenHash string     `json:"token_hash" gorm:"type:text;not null"`
	Type      string     `json:"type" gorm:"type:varchar(50);not null"` // Contoh: "bearer", "refresh", dll.
	ExpiredAt time.Time  `json:"expired_at" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (PersonalAccessToken) TableName() string {
	return "personal_access_tokens"
}

// Hook GORM untuk otomatis generate UUID pada PersonalAccessToken sebelum disimpan
func (t *PersonalAccessToken) BeforeCreate(tx *gorm.DB) (err error) {
	if uuid.UUID(t.ID) == uuid.Nil {
		t.ID = UUIDBinary(uuid.New())
	}
	return
}