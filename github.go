package latest

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/hashicorp/go-version"
)

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
func (g *GithubTag) Fetch() ([]*version.Version, []string, error) {

	var versions []*version.Version
	var malformedTags []string

	// Create a client
	client := g.newClient()
	tags, resp, err := client.Repositories.ListTags(g.Owner, g.Repository, nil)
	if err != nil {
		return versions, malformedTags, err
	}

	if resp.StatusCode != 200 {
		return versions, malformedTags, fmt.Errorf("Unknown status: %d", resp.StatusCode)
	}

	// fixF is FixVersionStrFunc transform tag name string into SemVer string
	// By default, it does nothing.
	fixF := g.fixVersionStrFunc()

	for _, tag := range tags {
		v, err := version.NewVersion(fixF(*tag.Name))
		if err != nil {
			malformedTags = append(malformedTags, fixF(*tag.Name))
			continue
		}
		versions = append(versions, v)
	}

	return versions, malformedTags, nil
}
