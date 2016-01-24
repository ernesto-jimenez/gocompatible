package importers

import (
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoDoc(t *testing.T) {
	tests := []struct {
		pkg            string
		recursive      bool
		expectsError   bool
		expectedResult []string
	}{
		{
			pkg:            "github.com/stretchr/testify",
			recursive:      false,
			expectsError:   false,
			expectedResult: []string{},
		},
		{
			pkg:          "github.com/stretchr/assert",
			recursive:    false,
			expectsError: false,
			expectedResult: []string{
				"bitbucket.org/pcas/math/testing/assert",
				"github.com/Billups/testify",
			},
		},
		{
			pkg:          "github.com/stretchr/testify",
			recursive:    true,
			expectsError: false,
			expectedResult: []string{
				"bitbucket.org/pcas/math/testing/assert",
				"github.com/Billups/testify",
				"github.com/Clever/optimus/tests",
				"github.com/EluctariLLC/dropbox",
				"github.com/Fiery/testify",
			},
		},
		{
			pkg:          "nonexistent",
			expectsError: true,
		},
	}

	s := httptest.NewServer(http.HandlerFunc(handleTestRequest))
	godoc := &GoDoc{}
	godoc.URL = s.URL
	defer s.Close()
	for _, test := range tests {
		l, err := godoc.List(test.pkg, test.recursive)
		if test.expectsError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, test.expectedResult, l)
		}
	}
}

func handleTestRequest(rw http.ResponseWriter, req *http.Request) {
	file := path.Base(req.URL.Path) + ".html"
	if strings.HasPrefix(req.URL.Path, "/importers/") {
		file = path.Base(req.URL.Path) + "-importers.json"
	}
	http.ServeFile(rw, req, path.Join("testdata", file))
}
