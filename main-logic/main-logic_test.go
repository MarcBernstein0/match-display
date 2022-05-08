package mainlogic

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	server *httptest.Server
	fetch  Fetch
)

func TestMain(m *testing.M) {
	fmt.Println("Starting mock server")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{status: UP}")
	}))

	fetch = New(server.URL, http.DefaultClient, time.Second)
	fmt.Println("Mock Server Running, Start Tests")
	m.Run()
}

func TestFetchData(t *testing.T) {
	tt := []struct {
		testName  string
		urlPath   string
		params    map[string]string
		wantData  string
		wantErr   error
		expectErr bool
	}{
		{
			testName:  "default test",
			urlPath:   "/test",
			params:    nil,
			wantData:  "{status: UP}",
			wantErr:   nil,
			expectErr: false,
		},
	}

	for _, test := range tt {
		t.Run(test.testName, func(t *testing.T) {
			t.Parallel()
			gotData, gotErr := fetch.FetchData(context.Background(), test.urlPath, test.params)
			assert.Equal(t, test.wantData, gotData)
			// assert.Equal(t, test.wantData, gotData)
			if test.expectErr {
				assert.EqualError(t, gotErr, test.wantErr.Error(), "expected %v but got %v", test.wantErr.Error(), gotErr.Error())
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}
