package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shaunpua/updoc/internal/doc"
)

type OrganizationService struct {
	orgRepo  doc.OrganizationRepository
	userRepo doc.UserRepository
}

func NewOrganizationService(orgRepo doc.OrganizationRepository, userRepo doc.UserRepository) *OrganizationService {
	return &OrganizationService{
		orgRepo:  orgRepo,
		userRepo: userRepo,
	}
}

type CreateOrgRequest struct {
	Name      string `json:"name" validate:"required,min=2,max=100"`
	UserName  string `json:"user_name" validate:"required,min=2,max=100"`
	UserEmail string `json:"user_email" validate:"required,email"`
	
	// Optional Confluence Integration
	ConfluenceBaseURL  string `json:"confluence_base_url,omitempty"`
	ConfluenceEmail    string `json:"confluence_email,omitempty"`
	ConfluenceToken    string `json:"confluence_token,omitempty"`
	ConfluenceSpaceKey string `json:"confluence_space_key,omitempty"`
}

type CreateOrgResponse struct {
	Organization *doc.Organization `json:"organization"`
	User         *doc.User         `json:"user"`
}

func (s *OrganizationService) CreateWithUser(ctx context.Context, req CreateOrgRequest) (*CreateOrgResponse, error) {
	// Generate slug from organization name
	slug := generateSlug(req.Name)

	// Check if organization already exists
	existing, err := s.orgRepo.GetBySlug(slug)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("organization with slug '%s' already exists", slug)
	}

	// Create organization
	org := &doc.Organization{
		Name:               req.Name,
		Slug:               slug,
		CreatedAt:          time.Now(),
		ConfluenceBaseURL:  req.ConfluenceBaseURL,
		ConfluenceEmail:    req.ConfluenceEmail,
		ConfluenceToken:    req.ConfluenceToken,
		ConfluenceSpaceKey: req.ConfluenceSpaceKey,
	}

	if err := s.orgRepo.Create(org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// Create admin user
	user := &doc.User{
		Email:     req.UserEmail,
		Name:      req.UserName,
		OrgID:     org.ID,
		Role:      "admin",
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &CreateOrgResponse{
		Organization: org,
		User:         user,
	}, nil
}

func (s *OrganizationService) GetBySlug(ctx context.Context, slug string) (*doc.Organization, error) {
	org, err := s.orgRepo.GetBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}
	return org, nil
}

func (s *OrganizationService) AddUser(ctx context.Context, orgID, email, name, role string) (*doc.User, error) {
	// Check if organization exists
	org, err := s.orgRepo.GetByID(orgID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	// Check if user already exists
	existing, err := s.userRepo.GetByEmail(email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("user with email '%s' already exists", email)
	}

	// Create user
	user := &doc.User{
		Email:     email,
		Name:      name,
		OrgID:     org.ID,
		Role:      role,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Helper function to generate URL-friendly slug
func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "&", "and")
	// Remove special characters (keep only alphanumeric and hyphens)
	result := ""
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result += string(char)
		}
	}
	return result
}
