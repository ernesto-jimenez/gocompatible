package util

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

// StripGopath takes a full path
// and removes the GOPATH to get
// the package's import path
func StripGopath(p string) string {
	gopaths := strings.Split(os.Getenv("GOPATH"), ":")
	for _, gopath := range gopaths {
		gopath, _ = filepath.Abs(gopath)
		gopath = path.Join(gopath, "src")
		if stripped := strings.TrimPrefix(p, gopath); stripped != p {
			return stripped
		}
	}
	return p
}
