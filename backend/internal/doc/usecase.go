package doc

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/shaunpua/updoc/internal/providers/confluence"
)

var defaultUser = User{ID: "u-123", Name: "demo-user"}

type Service struct {
	Client *resty.Client
	Store  FlagStore
}

func New(client *resty.Client, store FlagStore) *Service {
	return &Service{Client: client, Store: store}
}

func (s *Service) GetPage(id string) (*confluence.PageResp, error) {
	return confluence.GetPage(s.Client, id)
}

func (s *Service) UpdatePage(pageID, html, status string) (DocFlag, error) {
	if html != "" {
		p, err := confluence.GetPage(s.Client, pageID)
		if err != nil {
			return DocFlag{}, err
		}
		if err := confluence.UpdateBody(s.Client, p, html); err != nil {
			return DocFlag{}, err
		}
	}

	flag := DocFlag{
		ID:        uuid.NewString(),
		PageID:    pageID,
		Status:    status,
		Body:      html,
		User:      defaultUser,
		UpdatedAt: time.Now(),
	}
	s.Store.Save(flag)
	return flag, nil
}

func (s *Service) ListFlags(pageID string) ([]DocFlag, error) {
	return s.Store.ByPage(pageID)
}
