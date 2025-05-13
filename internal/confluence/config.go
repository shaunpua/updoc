// internal/confluence/config.go
package confluence

import "github.com/go-resty/resty/v2"

type Space struct {
	PageID string
	Client *resty.Client
}

func NewSpace(pageID string) *Space {
	return &Space{
		PageID: pageID,
		Client: New(), // reuse the HTTP client factory you wrote
	}
}
