package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          UUIDBinary   `json:"id" gorm:"type:binary(16);primaryKey"`
	Name        string       `json:"name" gorm:"unique;not null"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;foreignKey:ID;joinForeignKey:role_id;References:ID;joinReferences:permission_id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (Role) TableName() string {
	return "roles"
}

// Hook GORM untuk otomatis generate UUID pada Role
func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if uuid.UUID(r.ID) == uuid.Nil {
		r.ID = UUIDBinary(uuid.New())
	}
	return
}

type Permission struct {
	ID        UUIDBinary `json:"id" gorm:"type:binary(16);primaryKey"`
	Name      string     `json:"name" gorm:"unique;not null"`
	Roles     []Role     `json:"roles" gorm:"many2many:role_permissions;foreignKey:ID;joinForeignKey:permission_id;References:ID;joinReferences:role_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (Permission) TableName() string {
	return "permissions"
}

// Hook GORM untuk otomatis generate UUID pada Permission
func (p *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	if uuid.UUID(p.ID) == uuid.Nil {
		p.ID = UUIDBinary(uuid.New())
	}
	return
}