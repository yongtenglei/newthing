package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type TokenSession struct {
	gorm.Model

	Uuid         uuid.UUID `gorm:"index;not null" json:"uuid,required"`
	Mobile       string    `gorm:"type:varchar(32);index;not null" json:"mobile,required"`
	RefreshToken string    `gorm:"type:text;index;not null"`
	Issuer       string    `gorm:"type:varchar(32);not null" json:"issuer"`
	UserAgent    string    `gorm:"type:varchar(32);not null"`
	ClientIP     string    `gorm:"type:varchar(32);not null"`
	IssuedAt     time.Time `gorm:"default:now()"`
	ExpiredAt    int64     `gorm:"type:int;not null"`
}
