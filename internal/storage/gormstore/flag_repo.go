package gormstore

import (
	"github.com/google/uuid"
	"github.com/shaunpua/updoc/internal/doc"
	"gorm.io/gorm"
)

type FlagRepo struct{ DB *gorm.DB }

func NewFlagRepo(db *gorm.DB) *FlagRepo { return &FlagRepo{DB: db} }

func (r *FlagRepo) Save(f doc.DocFlag) error {
	rec := Flag{
		FlagID:    uuid.NewString(),
		PageID:    f.PageID,
		Status:    f.Status,
		Body:      f.Body,
		UserID:    f.User.ID,
		UpdatedAt: f.UpdatedAt,
	}
	return r.DB.Create(&rec).Error
}

func (r *FlagRepo) ByPage(pageID string) ([]doc.DocFlag, error) {
	var rows []Flag
	if err := r.DB.Where("page_id = ?", pageID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]doc.DocFlag, len(rows))
	for i, r := range rows {
		out[i] = doc.DocFlag{
			ID:        r.FlagID,
			PageID:    r.PageID,
			Status:    r.Status,
			Body:      r.Body,
			User:      doc.User{ID: r.UserID}, // name can be joined later
			UpdatedAt: r.UpdatedAt,
		}
	}
	return out, nil
}
