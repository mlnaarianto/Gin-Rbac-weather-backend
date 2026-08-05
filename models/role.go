package models

import (
	"time"
)

type Role struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name" gorm:"unique;not null"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;foreignKey:ID;joinForeignKey:role_id;References:ID;joinReferences:permission_id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (Role) TableName() string {
	return "roles"
}

type Permission struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"unique;not null"`
	Roles     []Role    `json:"roles" gorm:"many2many:role_permissions;foreignKey:ID;joinForeignKey:permission_id;References:ID;joinReferences:role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Permission) TableName() string {
	return "permissions"
}