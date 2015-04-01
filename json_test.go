package latest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON_implement(t *testing.T) {
	var _ Source = &JSON{}
}

func TestJSON_Validate(t *testing.T) {
	j := &JSON{
		URL: "http://example.com/info",
	}

	err := j.Validate()
	if err != nil {
		t.Fatalf("expect %s to eq nil", err.Error())
	}
}

func TestJSON_Fetch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"version":"1.2.4"}`)
	}))
	defer ts.Close()

	j := JSON{
		URL: ts.URL,
	}

	versions, malformed, err := j.Fetch()
	if err != nil {
		t.Fatalf("expect %s to eq nil", err.Error())
	}

	if len(malformed) != 0 {
		t.Fatalf("expect %d to eq 0", len(malformed))
	}

	expect := "1.2.4"
	if versions[0].String() != expect {
		t.Fatalf("expect %s to eq %s", versions[0].String(), expect)
	}
}
