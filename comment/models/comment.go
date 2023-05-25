package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	ID       string `gorm:"primarykey"`
	Entity   string `gorm:"column=entity"`
	EntityID string `gorm:"column=entity_id"`
	Comment  string `gorm:"column=comment"`
	UserId   string `gorm:"column=user_id"`
}
