package latest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTML_implement(t *testing.T) {
	var _ Source = &HTML{}
}

func TestHTML_Validate(t *testing.T) {
	h := &HTML{
		URL: "http://example.com/info",
	}

	err := h.Validate()
	if err != nil {
		t.Fatalf("expect %s to eq nil", err.Error())
	}
}

func TestHTML_Fetch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHTML := `<html><meta></html>`
		fmt.Fprintf(w, testHTML)
	}))
	defer ts.Close()

	h := HTML{
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
