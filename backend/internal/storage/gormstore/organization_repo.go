package gormstore

import (
	"github.com/shaunpua/updoc/internal/doc"
	"gorm.io/gorm"
)

type OrganizationRepo struct {
	DB *gorm.DB
}

func NewOrganizationRepo(db *gorm.DB) *OrganizationRepo {
	return &OrganizationRepo{DB: db}
}

func (r *OrganizationRepo) Create(org *doc.Organization) error {
	dbOrg := Organization{
		Name:               org.Name,
		Slug:               org.Slug,
		ConfluenceBaseURL:  org.ConfluenceBaseURL,
		ConfluenceEmail:    org.ConfluenceEmail,
		ConfluenceToken:    org.ConfluenceToken,
		ConfluenceSpaceKey: org.ConfluenceSpaceKey,
	}
	
	if err := r.DB.Create(&dbOrg).Error; err != nil {
		return err
	}
	
	// Update the domain object with generated values
	org.ID = dbOrg.ID
	org.CreatedAt = dbOrg.CreatedAt
	return nil
}

func (r *OrganizationRepo) GetBySlug(slug string) (*doc.Organization, error) {
	var dbOrg Organization
	if err := r.DB.Where("slug = ?", slug).First(&dbOrg).Error; err != nil {
		return nil, err
	}
	return r.toDomain(dbOrg), nil
}

func (r *OrganizationRepo) GetByID(id string) (*doc.Organization, error) {
	var dbOrg Organization
	if err := r.DB.Where("id = ?", id).First(&dbOrg).Error; err != nil {
		return nil, err
	}
	return r.toDomain(dbOrg), nil
}

func (r *OrganizationRepo) toDomain(o Organization) *doc.Organization {
	return &doc.Organization{
		ID:                 o.ID,
		Name:               o.Name,
		Slug:               o.Slug,
		CreatedAt:          o.CreatedAt,
		ConfluenceBaseURL:  o.ConfluenceBaseURL,
		ConfluenceEmail:    o.ConfluenceEmail,
		ConfluenceToken:    o.ConfluenceToken,
		ConfluenceSpaceKey: o.ConfluenceSpaceKey,
	}
}
