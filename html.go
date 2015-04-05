package latest

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-version"
)

// HTML is implemented Source interface.
// It fetch version information fron HTML and Scrap.
type HTML struct {
	// URL is HTML page URL which include version information
	URL string

	// Scraper is HTMLScraper which scrap HTML and extract version information.
	// By default, it does nothing (call)
	Scraper HTMLScraper
}

// HTMLScraper is inferface to scrap HTML and extract version information.
type HTMLScraper interface {
	Exec(r io.Reader) ([]string, *Meta, error)
}

type DefaultHTMLScrap struct{}

func (s *DefaultHTMLScrap) Exec(r io.Reader) ([]string, *Meta, error) {
	meta := &Meta{}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return []string{}, meta, err
	}

	b = bytes.Replace(b, []byte("\n"), []byte(""), -1)
	return []string{string(b[:])}, meta, nil
}

func (h *HTML) scraper() HTMLScraper {
	if h.Scraper == nil {
		return &DefaultHTMLScrap{}
	}

	return h.Scraper
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

// Fetch fetches HTML page and scrap it by HTMLScraper and extract version infomation
func (h *HTML) Fetch() (*FetchResponse, error) {

	fr := NewFetchResponse()

	// URL is validated before call
	u, _ := url.Parse(h.URL)

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

	scraper := h.scraper()
	verStrs, meta, err := scraper.Exec(resp.Body)
	if err != nil {
		return fr, err
	}

	if len(verStrs) == 0 {
		return fr, fmt.Errorf("version info is not found on %s", h.URL)
	}

	for _, verStr := range verStrs {
		v, err := version.NewVersion(verStr)
		if err != nil {
			fr.Malformeds = append(fr.Malformeds, verStr)
			continue
		}
		fr.Versions = append(fr.Versions, v)
	}

	fr.Meta = meta

	return fr, nil
}
