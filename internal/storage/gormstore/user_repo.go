package gormstore

import "gorm.io/gorm"

type UserRepo struct{ DB *gorm.DB }

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{DB: db} }

func (r *UserRepo) Ensure(id, name string) error {
	return r.DB.
		FirstOrCreate(&User{UserID: id, Name: name}, "user_id = ?", id).
		Error
}
