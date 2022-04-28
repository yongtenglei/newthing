package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Mobile   string `gorm:"type:varchar(20);index" json:"mobile,required"`
	Password string `gorm:"type:varchar(16)" json:"password,required"`
	Name     string `gorm:"type:varchar(20)" json:"name,required"`
	Gender   int    `gorm:"type:int" json:"gender,required"` // 0 for female 1 for male
	Mail     string `gorm:"type:varchar(36)" json:"mail,omitempty"`
}
