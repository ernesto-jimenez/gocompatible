package tester

import (
	"io/ioutil"
	"os"
)

// NewTempTest returns a Test with Gopath set to a clean temporary folder.
// It returns an error when the temporary folder cannot be created.
func NewTempTest() (*Test, error) {
	dir, err := ioutil.TempDir(os.TempDir(), "gocompile")
	if err != nil {
		return nil, err
	}
	return &Test{Gopath: dir}, nil
}
