// Package tester provides the functionality needed to test a package.
package tester

import (
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Test contains the details to fetch packages and run tests
type Test struct {
	// Gopath used to fetch the package and run the tests
	Gopath string

	// Stdout writer to get the command's output
	Stdout io.Writer
	// Stderr writer to get the command's output
	Stderr io.Writer
}

func (t *Test) env() []string {
	env := os.Environ()
	result := make([]string, len(env))
	for _, v := range env {
		if strings.HasPrefix(v, "GOPATH=") {
			continue
		}
		if strings.HasPrefix(v, "GOROOT=") {
			continue
		}
		result = append(result, v)
	}
	result = append(result, "GOPATH="+t.Gopath)
	return result
}

// Get fetches the package and its dependencies
func (t *Test) Get(pkg string) error {
	cmd := exec.Command("go", "get", "-d", "-t", pkg)
	t.prepareCmd(cmd)
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("go", "get", "-t", "./...")
	cmd.Dir = path.Join(t.Gopath, "src", pkg)
	t.prepareCmd(cmd)
	return cmd.Run()
}

// Test runs the tests form the package
func (t *Test) Test(pkg string) error {
	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = path.Join(t.Gopath, "src", pkg)
	t.prepareCmd(cmd)
	return cmd.Run()
}

// Build builds the package
func (t *Test) Build(pkg string) error {
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = path.Join(t.Gopath, "src", pkg)
	t.prepareCmd(cmd)
	return cmd.Run()
}

func (t *Test) prepareCmd(cmd *exec.Cmd) {
	cmd.Env = t.env()
	if t.Stdout != nil {
		cmd.Stdout = t.Stdout
	}
	if t.Stderr != nil {
		cmd.Stderr = t.Stderr
	}
}
