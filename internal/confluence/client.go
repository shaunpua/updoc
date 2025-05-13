package confluence

import (
	"os"

	"github.com/go-resty/resty/v2"
)

func New() *resty.Client {
	return resty.New().
		SetBaseURL(os.Getenv("CONF_BASE")). // https://updocu.atlassian.net/wiki
		SetBasicAuth(os.Getenv("CONF_EMAIL"),
			os.Getenv("CONF_TOKEN")).
		SetHeader("Accept", "application/json")
}
