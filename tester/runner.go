package tester

import (
	"errors"
	"fmt"
	"io"
)

// RunResult contains the results from using Run
type RunResult struct {
	OK          []string
	FailedGet   []string
	FailedBuild []string
	FailedTests []string
}

// TotalFailed returns the amount of packages that failed testing
func (r RunResult) TotalFailed() int {
	return len(r.FailedGet) + len(r.FailedBuild) + len(r.FailedTests)
}

// Total returns the amount of packages that were tested
func (r RunResult) Total() int {
	return r.TotalFailed() + len(r.OK)
}

// Err returns and error if any package failed testing. Nil otherwise.
func (r RunResult) Err() error {
	if r.TotalFailed() > 0 {
		return errors.New(r.String())
	}
	return nil
}

// String prints a one liner summary of the run
func (r RunResult) String() string {
	total := r.Total()
	failed := r.TotalFailed()
	if failed == 0 {
		return fmt.Sprintf("%d packages %d failed", total, failed)
	}
	return fmt.Sprintf(
		"%d packages %d failed: %4d failed get %4d build %4d test",
		total, failed,
		len(r.FailedGet), len(r.FailedBuild), len(r.FailedTests),
	)
}

// Run takes a list of packages and uses the given
// Test to fetch and run all the packages tests.
// It prints progress to the given io.Writer.
func Run(t *Test, pkgs []string, w io.Writer) *RunResult {
	res := &RunResult{}
	for _, pkg := range pkgs {
		if err := t.Get(pkg); err != nil {
			res.FailedGet = append(res.FailedGet, pkg)
			fmt.Fprintf(w, "%-7s %s - %s: %s\n", "FAIL", pkg, "go get", err)
			continue
		}
		if err := t.Build(pkg); err != nil {
			res.FailedBuild = append(res.FailedBuild, pkg)
			fmt.Fprintf(w, "%-7s %s - %s: %s\n", "FAIL", pkg, "go build", err)
			continue
		}
		if err := t.Test(pkg); err != nil {
			res.FailedTests = append(res.FailedTests, pkg)
			fmt.Fprintf(w, "%-7s %s - %s: %s\n", "FAIL", pkg, "go test", err)
			continue
		}
		res.OK = append(res.OK, pkg)
		fmt.Fprintf(w, "%-7s %s\n", "ok", pkg)
	}
	return res
}
