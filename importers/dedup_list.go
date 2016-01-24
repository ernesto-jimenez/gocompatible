package importers

import (
	"sort"
	"strings"
)

type deduplist map[string]bool

type sorted []string

func (l sorted) Len() int           { return len(l) }
func (l sorted) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l sorted) Less(i, j int) bool { return l[i] < l[j] }

func (d deduplist) add(path string) {
	d[path] = true
}

func (d deduplist) list() []string {
	var list []string
	for path := range d {
		list = append(list, path)
	}
	sort.Sort(sorted(list))
	last := "!"
	filtered := []string{}
	for _, pkg := range list {
		if strings.HasPrefix(pkg, last) {
			continue
		}
		filtered = append(filtered, pkg)
		last = pkg
	}
	return filtered
}
