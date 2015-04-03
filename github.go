package latest

import (
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"github.com/hashicorp/go-version"
)

// FixVersionStrFunc is function to fix version string
// so that it can be interpreted as SemVer by hashicorp/go-version
type FixVersionStrFunc func(string) string

var defaultFixVersionStrFunc FixVersionStrFunc

func init() {
	defaultFixVersionStrFunc = FixNothing()
}

// GithubTag is implemented Source interface. It uses GitHub API
// and fetch tags from repository.
type GithubTag struct {
	// Owner and Repository are GitHub owner name and its repository name
	// e.g., If you want to check https://github.com/tcnksm/ghr version
	// Repository is `ghr`, and Owner is `tcnksm`
	Owner      string
	Repository string

	// FixVersionStrFunc is function to fix version string (in this case tag
	// name string) on GitHub so that it can be interpreted as SemVer
	// by hashicorp/go-version. By default, it does nothing (calles FixNothing()).
	FixVersionStrFunc FixVersionStrFunc

	// URL & Token is used for GitHub Enterprise
	URL   string
	Token string
}

func (g *GithubTag) fixVersionStrFunc() FixVersionStrFunc {
	if g.FixVersionStrFunc == nil {
		return defaultFixVersionStrFunc
	}

	return g.FixVersionStrFunc
}

func FixNothing() FixVersionStrFunc {
	return func(s string) string {
		return s
	}
}

// DeleteFrontV delete first `v` charactor on version string
// e.g., `v0.1.1` becomes `0.1.1`
func DeleteFrontV() FixVersionStrFunc {
	return func(s string) string {
		return strings.Replace(s, "v", "", 1)
	}
}

func (g *GithubTag) newClient() *github.Client {
	return github.NewClient(nil)
}

func (g *GithubTag) Validate() error {

	if len(g.Repository) == 0 {
		return fmt.Errorf("GitHub repository name must be set")
	}

	if len(g.Owner) == 0 {
		return fmt.Errorf("GitHub owner name must be set")
	}

	return nil
}

// Fetch fetches github tags and interpret them as version.Version and return.
// To fetch tags, use google/go-github package.
func (g *GithubTag) Fetch() (*FetchResponse, error) {

	fr := NewFetchResponse()

	// Create a client
	client := g.newClient()
	tags, resp, err := client.Repositories.ListTags(g.Owner, g.Repository, nil)
	if err != nil {
		return fr, err
	}

	if resp.StatusCode != 200 {
		return fr, fmt.Errorf("Unknown status: %d", resp.StatusCode)
	}

	// fixF is FixVersionStrFunc transform tag name string into SemVer string
	// By default, it does nothing.
	fixF := g.fixVersionStrFunc()

	for _, tag := range tags {
		v, err := version.NewVersion(fixF(*tag.Name))
		if err != nil {
			fr.Malformeds = append(fr.Malformeds, fixF(*tag.Name))
			continue
		}
		fr.Versions = append(fr.Versions, v)
	}

	return fr, nil
}
