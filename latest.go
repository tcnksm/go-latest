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
// so that it can be parsed as Semantic version by hashicorp/go-version
type FixVersionStrFunc func(string) string

var defaultFixVersionStrFunc FixVersionStrFunc

func init() {
	// Doing nothing by default
	defaultFixVersionStrFunc = func(s string) string { return s }
}

type Source interface {
	// Validate validates Source option values
	// e.g., mandatory variables are set or not
	Validate() error

	// Fetch fetches version information from its source
	// and convert it to []*version.Version
	Fetch() ([]*version.Version, []string, error)
}

// CheckResponse stores check result
type CheckResponse struct {
	// Current is current version or tag on source
	Current string

	// Latest is true when target is greater than current on
	// source.
	Latest bool

	// New is true when target is greater than current on
	// source and new (not exist).
	New bool

	// Malformed store versions or tags which can not be parsed as
	// Semantic versioning (not compared with target)
	Malformed []string
}

// CheckLatest fetches last version information from its source
// And comapre with target and return results
func Check(target string, s Source) (*CheckResponse, error) {
	// Convert target to *version.Version
	targetV, err := version.NewVersion(target)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s : %s", err.Error())
	}

	// Validate options
	if err = s.Validate(); err != nil {
		return nil, err
	}

	versions, malformed, err := s.Fetch()
	if err != nil {
		return nil, err
	}

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

// DeleteFrontV delete first `v` charactor on version string
//
// e.g., `v0.1.1` becomes `0.1.1`
func DeleteFrontV() FixVersionStrFunc {
	return func(s string) string {
		return strings.Replace(s, "v", "", 1)
	}
}
