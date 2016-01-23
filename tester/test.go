// Package tester provides the functionality needed to test a package.
package tester

import (
	"os"
	"os/exec"
	"strings"
)

// Test contains the details to fetch packages and run tests
type Test struct {
	// Gopath used to fetch the package and run the tests
	Gopath string
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
	cmd := exec.Command("go", "get", pkg)
	cmd.Env = t.env()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Test runs the tests form the package
func (t *Test) Test(pkg string) error {
	cmd := exec.Command("go", "test", pkg)
	cmd.Env = t.env()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Build builds the package
func (t *Test) Build(pkg string) error {
	cmd := exec.Command("go", "build", pkg)
	cmd.Env = t.env()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
