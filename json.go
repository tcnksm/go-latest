package latest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-version"
)

// JSON is simple source which is assumes that
// a request returns json response which include `version` field.
type JSON struct {
	// URL is request URL to fetch version information
	URL string

	// FixVersionStrFunc transforms version string
	// so that it can be persed as semantic versioning
	// by hashicorp/go-version
	FixVersionStrFunc FixVersionStrFunc
}

type Response struct {
	Version string `json:"version"`
}

func (j *JSON) fixVersionStrFunc() FixVersionStrFunc {
	if j.FixVersionStrFunc == nil {
		return defaultFixVersionStrFunc
	}

	return j.FixVersionStrFunc
}

func (j *JSON) Validate() error {

	if len(j.URL) == 0 {
		return fmt.Errorf("URL must be set")
	}

	// Check URL can be parsed
	if _, err := url.Parse(j.URL); err != nil {
		return fmt.Errorf("%s is invalid URL: %s", j.URL, err.Error())
	}

	return nil
}

func (j *JSON) Fetch() ([]*version.Version, []string, error) {

	var versions []*version.Version
	var malformed []string

	u, _ := url.Parse(j.URL)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return versions, malformed, err
	}

	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return versions, malformed, err
	}

	if resp.StatusCode != 200 {
		return versions, malformed, fmt.Errorf("Unknown status: %d", resp.StatusCode)
	}

	var result Response
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&result); err != nil {
		return versions, malformed, err
	}

	verStr := result.Version

	fixF := j.fixVersionStrFunc()

	v, err := version.NewVersion(fixF(verStr))
	if err != nil {
		malformed = append(malformed, verStr)
		return versions, malformed, err
	}

	versions = append(versions, v)
	return versions, malformed, nil
}
