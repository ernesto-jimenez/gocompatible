package importers

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/gophergala2016/gocompatible/util"
)

// Local dependents matcher
type Local struct {
	// Path to dependents in.
	// By default it will search in all $GOPATH src folders
	Path string
}

// List looks for all dependent packages in the specified Path
func (g *Local) List(pkg string, recursive bool) ([]string, error) {
	searched := `\"` + pkg
	if !recursive {
		searched += `\"`
	}
	dirs := g.dirs()
	if len(dirs) == 0 {
		return nil, fmt.Errorf("no packages within scope %s", g.Path)
	}
	cmd := exec.Command("/bin/sh", "-c",
		fmt.Sprintf(
			"find %s -type f -name \\*.go -exec grep -F %s -sl {} +",
			strings.Join(dirs, " "),
			searched,
		),
	)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s - %s", err, stderr.String())
	}
	scanner := bufio.NewScanner(&stdout)
	list := deduplist(make(map[string]bool))
	for scanner.Scan() {
		p := util.StripGopath(path.Dir(scanner.Text()))
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
