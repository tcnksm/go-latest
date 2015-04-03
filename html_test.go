package latest

import (
	"io"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/hashicorp/go-version"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestHTML_implement(t *testing.T) {
	var _ Source = &HTML{}
}

// originalScrap is scrapFunc for test-fixtures/original.html
// It extracts VERSION from `<div class="version">VERSION</div>`
func originalScrap(r io.Reader) []string {

	// Check function attrs has correct class="version" key&value
	isTarget := func(targetVal string, attrs []html.Attribute) bool {
		for _, a := range attrs {
			if a.Namespace != "" {
				continue
			}

			if a.Key == "class" && a.Val == targetVal {
				return true
			}
		}
		return false
	}

	var verStrs []string
	z := html.NewTokenizer(r)

	for {
		switch z.Next() {
		case html.ErrorToken:
			return verStrs
		case html.StartTagToken:
			tok := z.Token()
			if tok.DataAtom == atom.Div && isTarget("version", tok.Attr) {
				z.Next()
				newTok := z.Token()
				verStrs = append(verStrs, newTok.String())
			}
		}
	}
}

func TestHTMLFetch(t *testing.T) {
	tests := []struct {
		testServer    *httptest.Server
		expectCurrent string
		scrapFunc     ScrapFunc
	}{
		{
			testServer:    fakeServer("test-fixtures/default.html"),
			expectCurrent: "1.2.3",
		},
		{
			testServer:    fakeServer("test-fixtures/original.html"),
			expectCurrent: "1.2.5",
			scrapFunc:     originalScrap,
		},
	}

	for i, tt := range tests {
		ts := tt.testServer
		defer ts.Close()

		h := &HTML{
			URL:       ts.URL,
			ScrapFunc: tt.scrapFunc,
		}

		fr, err := h.Fetch()
		if err != nil {
			t.Fatalf("#%d Fetch() expects error:%q to be nil", i, err.Error())
		}

		versions := fr.Versions
		if len(versions) == 0 {
			t.Fatalf("#%d Fetch() expects number of versions found from HTML not to be 0", i)
		}

		sort.Sort(version.Collection(versions))
		current := versions[len(versions)-1].String()
		if current != tt.expectCurrent {
			t.Fatalf("#%d Fetch() expects %s to be %s", i, current, tt.expectCurrent)
		}
	}

}
