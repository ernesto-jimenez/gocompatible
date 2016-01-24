package tester

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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
	tester.Stdout = os.Stdout
	tester.Stderr = os.Stderr

	require.NoError(t, tester.Get(pkg), "should get package")
	require.NoError(t, tester.Build(pkg), "package should build")
	require.NoError(t, tester.Test(pkg), "should pass tests")
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
