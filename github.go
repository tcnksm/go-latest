package latest

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/hashicorp/go-version"
)

// GithubTag store values related to Github
type GithubTag struct {
	// Repository is GitHub repository name
	Repository string

	// Owner is GitHub repository owner name
	Owner string

	// FixVersionStrFunc transforms version string
	// so that it can be persed as semantic versioning
	// by hashicorp/go-version
	FixVersionStrFunc FixVersionStrFunc

	// URL & Token is used for GitHub Enterprise
	// But not implemeted yet...
	URL   string
	Token string
}

func (g *GithubTag) fixVersionStrFunc() FixVersionStrFunc {
	if g.FixVersionStrFunc == nil {
		return defaultFixVersionStrFunc
	}

	return g.FixVersionStrFunc
}

// newClient create client for sending reuqest to Github
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

	// Add GHE validation

	return nil
}

func (g *GithubTag) Fetch() ([]*version.Version, []string, error) {

	var versions []*version.Version
	var malformedTags []string

	// Create client
	client := g.newClient()
	tags, resp, err := client.Repositories.ListTags(g.Owner, g.Repository, nil)
	if err != nil {
		return versions, malformedTags, err
	}

	if resp.StatusCode != 200 {
		return versions, malformedTags, fmt.Errorf("Unknown status: %d", resp.StatusCode)
	}

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
