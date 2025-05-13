package confluence

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// pageResp matches the JSON we get back from
// GET /rest/api/content/{id}?expand=body.storage,metadata.properties,version
type pageResp struct {
	ID    string `json:"id"`
	Title string `json:"title"`

	Body struct {
		Storage struct {
			Value string `json:"value"` // HTML in “storage” representation
		} `json:"storage"`
	} `json:"body"`

	Version struct {
		Number int `json:"number"`
	} `json:"version"`

	Metadata struct {
		Properties map[string]struct {
			Value any `json:"value"`
		} `json:"properties"`
	} `json:"metadata"`
}

// GetPage fetches title, body (HTML) and properties.
// It returns a parsed pageResp or a descriptive error.
func GetPage(c *resty.Client, id string) (*pageResp, error) {
	resp, err := c.R().
		SetResult(&pageResp{}).
		SetError(&map[string]any{}). // capture JSON error body
		SetQueryParam("expand", "body.storage,metadata.properties,version").
		Get(fmt.Sprintf("/rest/api/content/%s", id))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("confluence %d: %v", resp.StatusCode(), resp.Error())
	}
	return resp.Result().(*pageResp), nil
}

// UpdateBody replaces the page body (HTML) and bumps the version by +1.
// Title is kept as-is; call UpdateTitle if you want to change it too.
func UpdateBody(c *resty.Client, p *pageResp, newHTML string) error {
	payload := map[string]any{
		"id":    p.ID,
		"type":  "page",
		"title": p.Title, // keep current title
		"body": map[string]any{
			"storage": map[string]any{
				"value":          newHTML,
				"representation": "storage",
			},
		},
		"version": map[string]any{
			"number": p.Version.Number + 1,
		},
	}

	_, err := c.R().
		SetBody(payload).
		Put(fmt.Sprintf("/rest/api/content/%s", p.ID))
	return err
}
