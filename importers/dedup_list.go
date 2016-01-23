package importers

import (
	"sort"
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
	list := []string{}
	for path := range d {
		list = append(list, path)
	}
	sort.Sort(sorted(list))
	return list
}
