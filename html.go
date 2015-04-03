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

type HTML struct {
	URL       string
	ScrapFunc ScrapFunc
}

type ScrapFunc func(r io.Reader) ([]string, *Meta, error)

func ScrapNothing() ScrapFunc {
	return func(r io.Reader) ([]string, *Meta, error) {
		meta := &Meta{}
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return []string{}, meta, err
		}

		b = bytes.Replace(b, []byte("\n"), []byte(""), -1)
		return []string{string(b[:])}, meta, nil
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

	scrapFunc := h.scrapFunc()
	verStrs, meta, err := scrapFunc(resp.Body)
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
