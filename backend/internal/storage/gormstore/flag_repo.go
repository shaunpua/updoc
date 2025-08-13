package gormstore

import (
	"github.com/shaunpua/updoc/internal/doc"
	"gorm.io/gorm"
)

type FlagRepo struct{ DB *gorm.DB }

func NewFlagRepo(db *gorm.DB) *FlagRepo { return &FlagRepo{DB: db} }

// Implement new FlagRepository interface
func (r *FlagRepo) Create(flag *doc.Flag) error {
	dbFlag := Flag{
		DocumentID:  flag.DocumentID,
		CreatedBy:   flag.CreatedBy,
		AssignedTo:  flag.AssignedTo,
		Title:       flag.Title,
		Description: flag.Description,
		Priority:    flag.Priority,
		Status:      flag.Status,
		Resolution:  flag.Resolution,
		ResolvedAt:  flag.ResolvedAt,
		CreatedAt:   flag.CreatedAt,
		UpdatedAt:   flag.UpdatedAt,
	}

	if err := r.DB.Create(&dbFlag).Error; err != nil {
		return err
	}

	// Update the flag with the generated ID
	flag.ID = dbFlag.ID
	return nil
}

func (r *FlagRepo) GetByID(id string) (*doc.Flag, error) {
	var dbFlag Flag
	if err := r.DB.Preload("Creator").Preload("Assignee").Preload("Document").
		First(&dbFlag, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return r.toDomainFlag(dbFlag), nil
}

func (r *FlagRepo) GetByDocumentID(documentID string) ([]*doc.Flag, error) {
	var dbFlags []Flag
	if err := r.DB.Preload("Creator").Preload("Assignee").
		Where("document_id = ?", documentID).Find(&dbFlags).Error; err != nil {
		return nil, err
	}

	flags := make([]*doc.Flag, len(dbFlags))
	for i, dbFlag := range dbFlags {
		flags[i] = r.toDomainFlag(dbFlag)
	}
	return flags, nil
}

func (r *FlagRepo) GetByFilters(filters doc.FlagFilters) ([]*doc.Flag, error) {
	query := r.DB.Preload("Creator").Preload("Assignee").Preload("Document")

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Priority != "" {
		query = query.Where("priority = ?", filters.Priority)
	}
	if filters.AssignedTo != "" {
		query = query.Where("assigned_to = ?", filters.AssignedTo)
	}
	if filters.CreatedBy != "" {
		query = query.Where("created_by = ?", filters.CreatedBy)
	}
	if filters.WorkspaceID != "" {
		// Join with documents table to filter by workspace
		query = query.Joins("JOIN documents ON flags.document_id = documents.id").
			Where("documents.workspace_id = ?", filters.WorkspaceID)
	}
	if filters.Search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?",
			"%"+filters.Search+"%", "%"+filters.Search+"%")
	}

	var dbFlags []Flag
	if err := query.Find(&dbFlags).Error; err != nil {
		return nil, err
	}

	flags := make([]*doc.Flag, len(dbFlags))
	for i, dbFlag := range dbFlags {
		flags[i] = r.toDomainFlag(dbFlag)
	}
	return flags, nil
}

func (r *FlagRepo) Update(flag *doc.Flag) error {
	dbFlag := Flag{
		ID:          flag.ID,
		DocumentID:  flag.DocumentID,
		CreatedBy:   flag.CreatedBy,
		AssignedTo:  flag.AssignedTo,
		Title:       flag.Title,
		Description: flag.Description,
		Priority:    flag.Priority,
		Status:      flag.Status,
		Resolution:  flag.Resolution,
		ResolvedAt:  flag.ResolvedAt,
		UpdatedAt:   flag.UpdatedAt,
	}

	return r.DB.Save(&dbFlag).Error
}

// Helper method to convert GORM model to domain model
func (r *FlagRepo) toDomainFlag(dbFlag Flag) *doc.Flag {
	flag := &doc.Flag{
		ID:          dbFlag.ID,
		DocumentID:  dbFlag.DocumentID,
		CreatedBy:   dbFlag.CreatedBy,
		AssignedTo:  dbFlag.AssignedTo,
		Title:       dbFlag.Title,
		Description: dbFlag.Description,
		Priority:    dbFlag.Priority,
		Status:      dbFlag.Status,
		Resolution:  dbFlag.Resolution,
		ResolvedAt:  dbFlag.ResolvedAt,
		CreatedAt:   dbFlag.CreatedAt,
		UpdatedAt:   dbFlag.UpdatedAt,
	}

	// Convert related entities if loaded
	if dbFlag.Creator.ID != "" {
		flag.Creator = &doc.User{
			ID:    dbFlag.Creator.ID,
			Email: dbFlag.Creator.Email,
			Name:  dbFlag.Creator.Name,
			OrgID: dbFlag.Creator.OrgID,
			Role:  dbFlag.Creator.Role,
		}
	}

	if dbFlag.Assignee != nil && dbFlag.Assignee.ID != "" {
		flag.Assignee = &doc.User{
			ID:    dbFlag.Assignee.ID,
			Email: dbFlag.Assignee.Email,
			Name:  dbFlag.Assignee.Name,
			OrgID: dbFlag.Assignee.OrgID,
			Role:  dbFlag.Assignee.Role,
		}
	}

	if dbFlag.Document.ID != "" {
		flag.Document = &doc.Document{
			ID:          dbFlag.Document.ID,
			WorkspaceID: dbFlag.Document.WorkspaceID,
			Title:       dbFlag.Document.Title,
			URL:         dbFlag.Document.URL,
			ExternalID:  dbFlag.Document.ExternalID,
			OwnerID:     dbFlag.Document.OwnerID,
			LastChecked: dbFlag.Document.LastChecked,
			CreatedAt:   dbFlag.Document.CreatedAt,
		}
	}

	return flag
}
