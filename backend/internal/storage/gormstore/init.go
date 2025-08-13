package gormstore

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Organization{}, &User{}, &Workspace{}, &Document{}, &Flag{}, &Notification{})
}
