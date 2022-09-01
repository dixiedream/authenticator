package model

import (
	"time"

	"github.com/dixiedream/authenticator/utils"
	"gorm.io/gorm"
)

type Server struct {
	gorm.Model
	Name       string    `gorm:"unique; not null" json:"name" validate:"required"`
	Hostname   string    `json:"hostname" gorm:"index" validate:"required"`
	Password   string    `json:"password" validate:"required"`
	IpAddress  string    `gorm:"unique;index" json:"ipAddress" validate:"required"`
	LastAccess time.Time `json:"lastAccess"`
	Role       int       `gorm:"default:743" json:"role"`
	Sessions   []Session `json:"sessions"`
}

func (s *Server) BeforeCreate(tx *gorm.DB) error {
	password, err := utils.HashPassword(s.Password)
	if err != nil {
		return err
	}

	s.Password = password
	return nil
}
