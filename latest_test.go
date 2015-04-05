package latest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON_Check(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"version":"1.2.3"}`)
	}))
	defer ts.Close()

	j := &JSON{
		URL: ts.URL,
	}

	target := "1.2.0"
	res, err := Check(j, target)
	if err != nil {
		t.Fatalf("expect %q to be nil", err.Error())
	}

	if res.Latest {
		t.Fatalf("expect %t to be false", res.Latest)
	}

	expect := "1.2.3"
	if res.Current != expect {
		t.Fatalf("expect %q to be %q", res.Current, expect)
	}

}
