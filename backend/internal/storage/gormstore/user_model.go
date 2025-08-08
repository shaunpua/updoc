package gormstore

type User struct {
	UserID string `gorm:"primaryKey"`
	Name   string
	Flags  []Flag `gorm:"foreignKey:UserID"`
}
