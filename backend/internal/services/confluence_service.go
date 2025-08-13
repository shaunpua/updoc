package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/shaunpua/updoc/internal/doc"
)

type ConfluenceService struct {
	orgRepo doc.OrganizationRepository
}

func NewConfluenceService(orgRepo doc.OrganizationRepository) *ConfluenceService {
	return &ConfluenceService{orgRepo: orgRepo}
}

type ConfluenceTestResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// TestConnection tests the Confluence connection for an organization
func (s *ConfluenceService) TestConnection(ctx context.Context, orgID string) (*ConfluenceTestResponse, error) {
	org, err := s.orgRepo.GetByID(orgID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	if org.ConfluenceBaseURL == "" || org.ConfluenceEmail == "" || org.ConfluenceToken == "" {
		return &ConfluenceTestResponse{
			Success: false,
			Message: "Confluence integration not configured",
			Details: "Missing base URL, email, or token",
		}, nil
	}

	// Test connection by trying to get user info
	client := resty.New()
	resp, err := client.R().
		SetBasicAuth(org.ConfluenceEmail, org.ConfluenceToken).
		Get(org.ConfluenceBaseURL + "/rest/api/user/current")

	if err != nil {
		return &ConfluenceTestResponse{
			Success: false,
			Message: "Connection failed",
			Details: err.Error(),
		}, nil
	}

	if resp.StatusCode() != 200 {
		return &ConfluenceTestResponse{
			Success: false,
			Message: "Authentication failed",
			Details: fmt.Sprintf("HTTP %d: %s", resp.StatusCode(), resp.String()),
		}, nil
	}

	return &ConfluenceTestResponse{
		Success: true,
		Message: "Connection successful",
		Details: "Successfully authenticated with Confluence",
	}, nil
}

type ConfluencePageInfo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Space string `json:"space"`
}

// ListPages gets pages from the configured space
func (s *ConfluenceService) ListPages(ctx context.Context, orgID string, limit int) ([]ConfluencePageInfo, error) {
	org, err := s.orgRepo.GetByID(orgID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	if org.ConfluenceBaseURL == "" || org.ConfluenceEmail == "" || org.ConfluenceToken == "" {
		return nil, fmt.Errorf("confluence integration not configured")
	}

	if limit <= 0 {
		limit = 10
	}

	client := resty.New()
	
	url := fmt.Sprintf("%s/rest/api/content", org.ConfluenceBaseURL)
	req := client.R().SetBasicAuth(org.ConfluenceEmail, org.ConfluenceToken)
	
	if org.ConfluenceSpaceKey != "" {
		req.SetQueryParam("spaceKey", org.ConfluenceSpaceKey)
	}
	
	resp, err := req.
		SetQueryParam("limit", fmt.Sprintf("%d", limit)).
		SetQueryParam("expand", "space").
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch pages: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("confluence API error: HTTP %d", resp.StatusCode())
	}

	// Simple response parsing (you might want to use a proper struct)
	var result struct {
		Results []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			Links struct {
				WebUI string `json:"webui"`
			} `json:"_links"`
			Space struct {
				Key string `json:"key"`
			} `json:"space"`
		} `json:"results"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	pages := make([]ConfluencePageInfo, len(result.Results))
	for i, page := range result.Results {
		pages[i] = ConfluencePageInfo{
			ID:    page.ID,
			Title: page.Title,
			URL:   org.ConfluenceBaseURL + page.Links.WebUI,
			Space: page.Space.Key,
		}
	}

	return pages, nil
}
