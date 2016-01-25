package tester

import (
	"os/exec"
	"path"
)

// Checkout allows you to checkout a specific revision of the given package
// into the Test's gopath
func (t *Test) Checkout(pkg, rev string) error {
	wd := path.Join(t.Gopath, "src", pkg)
	cmd := exec.Command("git", "checkout", "-B", "gocompatible-test", rev)
	cmd.Dir = wd
	t.prepareCmd(cmd)
	return cmd.Run()
}

// FetchRemote allows you to fetch commits from another remote source
func (t *Test) FetchRemote(pkg, name, url string) error {
	wd := path.Join(t.Gopath, "src", pkg)
	cmd := exec.Command("git", "remote", "add", "-f", name, url)
	cmd.Dir = wd
	t.prepareCmd(cmd)
	return cmd.Run()
}
