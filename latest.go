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
	"strings"

	"github.com/hashicorp/go-version"
)

// FixVersionStrFunc is function to fix version string
// so that it can be interpreted as SemVer by hashicorp/go-version
type FixVersionStrFunc func(string) string

var defaultFixVersionStrFunc FixVersionStrFunc

func init() {
	defaultFixVersionStrFunc = FixNothing()
}

// Source is version information source like GitHub or your server HTML.
type Source interface {
	// Validate validates Source struct e.g., mandatory variables are set
	Validate() error

	// Fetch fetches version information from its source
	// and convert it into version.Version
	Fetch() ([]*version.Version, []string, error)
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
	Malformed []string
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

	versions, malformed, err := s.Fetch()
	if err != nil {
		return nil, err
	}

	// Source must has at leaset one version information
	if len(versions) == 0 {
		return nil, fmt.Errorf("no version to compare")
	}
	sort.Sort(version.Collection(versions))
	currentV := versions[len(versions)-1]

	var latest, new bool
	// If target >= current, target is `lastest`
	if targetV.GreaterThan(currentV) {
		latest = true

		// If target > current, target is `new`
		if !targetV.Equal(currentV) {
			new = true
		}
	}

	return &CheckResponse{
		Current:   currentV.String(),
		Latest:    latest,
		New:       new,
		Malformed: malformed,
	}, nil
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
