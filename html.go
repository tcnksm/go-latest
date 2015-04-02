package latest

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"bytes"

	"github.com/hashicorp/go-version"
)

type HTML struct {
	URL       string
	ScrapFunc ScrapFunc
}

type ScrapFunc func(r io.Reader) string

func ScrapNothing() ScrapFunc {
	return func(r io.Reader) string {
		b, _ := ioutil.ReadAll(r)
		b = bytes.Replace(b, []byte("\n"), []byte(""), -1)
		return string(b[:])
	}
}

var defaultScrapFunc ScrapFunc

func init() {
	defaultScrapFunc = ScrapNothing()
}

func (h *HTML) scrapFunc() ScrapFunc {
	if h.ScrapFunc == nil {
		return defaultScrapFunc
	}

	return h.ScrapFunc
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

	// URL is validated before call
	u, _ := url.Parse(h.URL)

	// Create a new request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return versions, malformed, err
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
		return versions, malformed, err
	}

	if resp.StatusCode != 200 {
		return versions, malformed, fmt.Errorf("unknown status: %d", resp.StatusCode)
	}

	scrapFunc := h.scrapFunc()
	verStr := scrapFunc(resp.Body)
	if len(verStr) == 0 {
		return versions, malformed, fmt.Errorf("version info is not found on %s", h.URL)
	}

	v, err := version.NewVersion(verStr)
	if err != nil {
		malformed = append(malformed, verStr)
		return versions, malformed, err
	}

	versions = append(versions, v)
	return versions, malformed, nil
}
