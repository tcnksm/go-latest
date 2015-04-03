/*
For a full guide visit http://github.com/tcnksm/go-latest

  package main

  import (
      "github.com/tcnksm/go-latest"
  )

*/
package latest

import (
	"fmt"
	"sort"

	"github.com/hashicorp/go-version"
)

// Source is version information source like GitHub or your server HTML.
type Source interface {
	// Validate validates Source struct e.g., mandatory variables are set
	Validate() error

	// Fetch fetches version information from its source
	// and convert it into version.Version
	Fetch() (*FetchResponse, error)
}

// FetchResponse stores Fetch() results
type FetchResponse struct {
	Versions   []*version.Version
	Malformeds []string
	Meta       *Meta
}

// Meta is meta information from source
type Meta struct {
	Message string
	URL     string
}

// CheckResponse stores check results
type CheckResponse struct {
	// Current is current latest version on source
	Current string

	// Latest is true when target is greater than Current on source.
	Latest bool

	// New is true when target is greater than Current on
	// source and new (not exist).
	New bool

	// Malformed store versions or tags which can not be parsed as
	// Semantic versioning (not compared with target)
	Malformeds []string

	//
	Meta *Meta
}

// CheckLatest fetches last version information from its source
// and compares with target and return result (CheckResponse)
func Check(target string, s Source) (*CheckResponse, error) {
	// Convert target to *version.Version
	targetV, err := version.NewVersion(target)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s : %s", err.Error())
	}

	// Validate source
	if err = s.Validate(); err != nil {
		return nil, err
	}

	fr, err := s.Fetch()
	if err != nil {
		return nil, err
	}

	// Source must has at leaset one version information
	versions := fr.Versions
	if len(fr.Versions) == 0 {
		return nil, fmt.Errorf("no version to compare")
	}
	sort.Sort(version.Collection(versions))
	currentV := versions[len(versions)-1]

	var latest, new bool
	// If target > current, target is `latest` and `new`
	if targetV.GreaterThan(currentV) {
		new, latest = true, true
	}

	// If target = current, target is `latest`
	if targetV.Equal(currentV) {
		latest = true
	}

	return &CheckResponse{
		Current:    currentV.String(),
		Latest:     latest,
		New:        new,
		Malformeds: fr.Malformeds,
		Meta:       fr.Meta,
	}, nil
}

func NewFetchResponse() *FetchResponse {
	var versions []*version.Version
	var malformeds []string
	return &FetchResponse{
		Versions:   versions,
		Malformeds: malformeds,
		Meta:       &Meta{},
	}
}
