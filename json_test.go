package latest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func fakeServer(fixture string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open(fixture)
		if err != nil {
			// Should not reach here
			panic(err)
		}
		io.Copy(w, f)
	}))
}

func TestJSON_implement(t *testing.T) {
	var _ Source = &JSON{}
}

func TestJSONValidate(t *testing.T) {

	tests := []struct {
		JSON      *JSON
		expectErr bool
	}{
		{
			JSON: &JSON{
				URL: "http://good.com",
			},
			expectErr: false,
		},
		{
			JSON: &JSON{
				URL: "",
			},
			expectErr: true,
		},
	}

	for i, tt := range tests {
		j := tt.JSON
		err := j.Validate()
		if tt.expectErr == (err == nil) {
			t.Fatalf("#%d Validate() expects err == nil to eq %t", i, tt.expectErr)
		}
	}
}

// OriginalResponse implements Receiver and receives test-fixtures/original.json
type OriginalResponse struct {
	Name        string `json:"name"`
	VersionInfo string `json:"version_info"`
}

func (r *OriginalResponse) Version() string {
	verStr := strings.Replace(r.VersionInfo, "v", "", 1)
	return verStr
}

func TestJSONFetch(t *testing.T) {

	tests := []struct {
		testServer    *httptest.Server
		receiver      Receiver
		expectCurrent string
	}{
		{
			testServer:    fakeServer("test-fixtures/default.json"),
			expectCurrent: "1.2.3",
		},
		{
			testServer:    fakeServer("test-fixtures/original.json"),
			expectCurrent: "0.1.0",
			receiver:      &OriginalResponse{},
		},
	}

	for i, tt := range tests {
		ts := tt.testServer
		defer ts.Close()

		j := &JSON{
			URL:      ts.URL,
			Receiver: tt.receiver,
		}

		versions, _, err := j.Fetch()
		if err != nil {
			t.Fatalf("#%d Fetch() expects error:%q to be nil", i, err.Error())
		}

		current := versions[0].String()
		if current != tt.expectCurrent {
			t.Fatalf("#%d Fetch() expects %s to be %s", i, current, tt.expectCurrent)
		}
	}
}
