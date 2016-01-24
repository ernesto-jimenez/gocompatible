package tester

import (
	"os/exec"
	"path"
)

// Checkout allows you to checkout a specific
// revision of the given package into the Test's
// gopath
func (t *Test) Checkout(pkg, rev string) error {
	cmd := exec.Command("git", "checkout", rev)
	cmd.Dir = path.Join(t.Gopath, "src", pkg)
	t.prepareCmd(cmd)
	return cmd.Run()
}
