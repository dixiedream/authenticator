package model

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	ServerID uint   `json:"serverId"`
	Expired  bool   `json:"expired"`
	Token    string `gorm:"unique;index" json:"token"`
}
