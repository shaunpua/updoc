package gormstore

import (
	"github.com/shaunpua/updoc/internal/doc"
	"gorm.io/gorm"
)

type UserRepo struct{ DB *gorm.DB }

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{DB: db} }

// Implement new UserRepository interface
func (r *UserRepo) Create(user *doc.User) error {
	dbUser := User{
		Email:    user.Email,
		Name:     user.Name,
		OrgID:    user.OrgID,
		Role:     user.Role,
		IsActive: true,
	}

	if err := r.DB.Create(&dbUser).Error; err != nil {
		return err
	}

	// Update the user with the generated ID
	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt
	return nil
}

func (r *UserRepo) GetByEmail(email string) (*doc.User, error) {
	var dbUser User
	if err := r.DB.Where("email = ?", email).First(&dbUser).Error; err != nil {
		return nil, err
	}

	return r.toDomainUser(dbUser), nil
}

func (r *UserRepo) GetByOrgID(orgID string) ([]*doc.User, error) {
	var dbUsers []User
	if err := r.DB.Where("org_id = ? AND is_active = true", orgID).Find(&dbUsers).Error; err != nil {
		return nil, err
	}

	users := make([]*doc.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = r.toDomainUser(dbUser)
	}
	return users, nil
}

func (r *UserRepo) GetByID(id string) (*doc.User, error) {
	var dbUser User
	if err := r.DB.Where("id = ?", id).First(&dbUser).Error; err != nil {
		return nil, err
	}

	return r.toDomainUser(dbUser), nil
}

// Helper method to convert GORM model to domain model
func (r *UserRepo) toDomainUser(dbUser User) *doc.User {
	return &doc.User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		Name:      dbUser.Name,
		OrgID:     dbUser.OrgID,
		Role:      dbUser.Role,
		IsActive:  dbUser.IsActive,
		CreatedAt: dbUser.CreatedAt,
	}
}

// Legacy method (for backward compatibility)
func (r *UserRepo) Ensure(id, name string) error {
	// This method doesn't make sense with the new structure since we need email and orgID
	// For now, return an error to force migration to new methods
	return gorm.ErrInvalidData
}
