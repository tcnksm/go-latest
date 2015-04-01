package latest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTMLMeta_implement(t *testing.T) {
	var _ Source = &HTMLMeta{}
}

func TestHTMLMeta_Validate(t *testing.T) {
	h := &HTMLMeta{
		URL: "http://example.com/info",
	}

	err := h.Validate()
	if err != nil {
		t.Fatalf("expect %s to eq nil", err.Error())
	}
}

func TestHTMLMeta_Fetch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHTML := `<html><meta></html>`
		fmt.Fprintf(w, testHTML)
	}))
	defer ts.Close()

	h := HTMLMeta{
		URL: ts.URL,
	}

	versions, malformed, err := h.Fetch()
	if err != nil {
		t.Fatalf("expect %s to eq nil", err.Error())
	}

	if len(malformed) != 0 {
		t.Fatalf("expect %d to eq 0", len(malformed))
	}

	_ = versions
}
