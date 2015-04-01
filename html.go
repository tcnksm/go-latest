package latest

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/go-version"
)

type HTML struct {
	// URL is request URL to fetch version information
	URL string

	// FixVersionStrFunc transforms version string
	// so that it can be persed as semantic versioning
	// by hashicorp/go-version
	FixVersionStrFunc FixVersionStrFunc
}

func (h *HTML) fixVersionStrFunc() FixVersionStrFunc {
	if h.FixVersionStrFunc == nil {
		return defaultFixVersionStrFunc
	}

	return h.FixVersionStrFunc
}

func (h *HTML) Validate() error {

	if len(h.URL) == 0 {
		return fmt.Errorf("URL must be set")
	}

	// Check URL can be parsed
	if _, err := url.Parse(h.URL); err != nil {
		return fmt.Errorf("%s is invalid URL: %s", h.URL, err.Error())
	}

	return nil
}

func (h *HTML) Fetch() ([]*version.Version, []string, error) {

	var versions []*version.Version
	var malformed []string

	return versions, malformed, nil
}
