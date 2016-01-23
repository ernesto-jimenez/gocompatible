package tester

import (
	"flag"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

var gopath string

func TestFailsWithoutGopath(t *testing.T) {
	pkg := "github.com/stretchr/testify"
	tester := &Test{}

	require.Error(t, tester.Get(pkg))
	require.Error(t, tester.Test(pkg))
}

func TestPasses(t *testing.T) {
	pkg := "github.com/stretchr/testify"
	tester := &Test{Gopath: gopath}

	require.NoError(t, tester.Get(pkg))
	require.NoError(t, tester.Test(pkg))
}

func TestGetFailsWithUnknownPackage(t *testing.T) {
	tester := &Test{Gopath: gopath}

	require.Error(t, tester.Get("whatever"))
}

func TestMain(m *testing.M) {
	flag.Parse()
	dir, err := ioutil.TempDir(os.TempDir(), "gocompile-tests")
	if err != nil {
		panic(err)
	}
	defer func() {
		os.RemoveAll(dir)
	}()
	gopath = dir
	os.Exit(m.Run())
}
