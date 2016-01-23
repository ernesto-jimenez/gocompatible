package importers

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Local dependents matcher
type Local struct {
	// Path to dependents in.
	// By default it will search in all $GOPATH src folders
	Path string
}

// List looks for all dependent packages in the specified Path
func (g *Local) List(pkg string, recursive bool) ([]string, error) {
	searched := `"` + pkg
	if !recursive {
		searched += `"`
	}
	args := []string{
		searched,
		"-RIsl",             // recursive, ignore binaries, no errors, only files
		"--include", "*.go", // check only go source files
		"--exclude-dir", "Godeps", "--exclude-dir", "vendor", // skip vendored deps
	}
	dirs := g.dirs()
	if len(dirs) == 0 {
		return nil, fmt.Errorf("no packages within scope %s", g.Path)
	}
	args = append(args, dirs...)
	cmd := exec.Command("grep", args...)
	r, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(r)
	list := deduplist(make(map[string]bool))
	for scanner.Scan() {
		p := stripGopath(path.Dir(scanner.Text()))
		list.add(p)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return list.list(), nil
}

func (g *Local) dirs() []string {
	gopath := strings.Split(os.Getenv("GOPATH"), ":")
	if g.Path == "" {
		return gopath
	}
	result := gopath[:0]
	for _, gp := range gopath {
		p := path.Join(gp, "src", g.Path)
		if stat, err := os.Stat(p); err == nil && stat.IsDir() {
			result = append(result, p)
		}
	}
	return result
}

func stripGopath(path string) string {
	parts := strings.SplitN(path, "/src/", 2)
	return parts[len(parts)-1]
}
