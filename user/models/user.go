package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       string `gorm:"primarykey"`
	Name     string `gorm:"column=name"`
	Email    string `gorm:"column=email"`
	Password string `gorm:"column=password"`
}
