package importers

type deduplist map[string]bool

func (d deduplist) add(path string) {
	d[path] = true
}

func (d deduplist) list() []string {
	list := []string{}
	for path := range d {
		list = append(list, path)
	}
	return list
}
