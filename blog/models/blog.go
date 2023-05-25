package models

import "gorm.io/gorm"

type Blog struct {
	gorm.Model
	ID          string `gorm:"primarykey"`
	Title       string `gorm:"column=title"`
	Description string `gorm:"column=description"`
	Category    string `gorm:"column=category"`
	AuthorId    string `gorm:"column=authorId"`
	Tags        string `gorm:"column=tags"`
}
