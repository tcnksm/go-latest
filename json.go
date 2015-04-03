package latest

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-version"
)

var (
	defaultDialTimeout = 5 * time.Second
)

// JSON is implemented Source interface. It fetches version infomation
// from URL.
type JSON struct {
	// URL is URL which return json with version information.
	URL string

	// Receiver is Receiver interface to use original json response.
	// Reveiver should be defined json response struct and Version()
	// to return SemVer format version string.
	// By default, DefaultResponse is used.
	Receiver Receiver
}

type Receiver interface {
	// VersionInfo() returns version list.
	// It must be SemVer format. If response is not SemVer format,
	// transform it in this function.
	VersionInfo() ([]string, error)

	// MetaInfo() returns Meta information
	MetaInfo() (*Meta, error)
}

// DefaultResponse assumes response include `version` field and version
// is SemVer format. e.g., {"version":"1.2.3"}
type DefaultResponse struct {
	Version string `json:"version"`
	Message string `json:"message"`
	URL     string `json:"url"`
}

func (res *DefaultResponse) VersionInfo() ([]string, error) {
	return []string{res.Version}, nil
}

func (res *DefaultResponse) MetaInfo() (*Meta, error) {
	return &Meta{
		Message: res.Message,
		URL:     res.URL,
	}, nil
}

func (j *JSON) receiver() Receiver {
	if j.Receiver == nil {
		return &DefaultResponse{}
	}

	return j.Receiver
}

func (j *JSON) Validate() error {

	if len(j.URL) == 0 {
		return fmt.Errorf("URL must be set")
	}

	// Check URL can be parsed by net.URL
	if _, err := url.Parse(j.URL); err != nil {
		return fmt.Errorf("%s is invalid URL: %s", j.URL, err.Error())
	}

	return nil
}

// Fetch fetches Json from server and interpret them as version.Version and return.
func (j *JSON) Fetch() (*FetchResponse, error) {

	fr := NewFetchResponse()

	// URL is validated before call
	u, _ := url.Parse(j.URL)

	// Create a new request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return fr, err
	}
	req.Header.Add("Accept", "application/json")

	// Create client
	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: func(n, a string) (net.Conn, error) {
			return net.DialTimeout(n, a, defaultDialTimeout)
		},
	}

	client := &http.Client{
		Transport: t,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fr, err
	}

	if resp.StatusCode != 200 {
		return fr, fmt.Errorf("unknown status: %d", resp.StatusCode)
	}

	result := j.receiver()
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&result); err != nil {
		return fr, err
	}

	verStrs, err := result.VersionInfo()
	if err != nil {
		return fr, err
	}

	if len(verStrs) == 0 {
		return fr, fmt.Errorf("version info is not found on %s", j.URL)
	}

	for _, verStr := range verStrs {
		v, err := version.NewVersion(verStr)
		if err != nil {
			fr.Malformeds = append(fr.Malformeds, verStr)
		}
		fr.Versions = append(fr.Versions, v)
	}

	fr.Meta, err = result.MetaInfo()
	if err != nil {
		return fr, err
	}

	return fr, nil
}
