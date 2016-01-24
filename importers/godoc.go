// Package importers provides functionality to track down packages depending on
// certain package.
package importers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/webdevdata/webdevdata-tools/webdevdata"
	"golang.org/x/net/html"
)

// GoDoc is used to extract importers information from the godoc website.
// By default information is extracted from https://godoc.org
type GoDoc struct {
	// URL contains the url used to extract the information from.
	// If empty it'll use https://godoc.org
	URL string
	// Path containing the dependents we want to list.
	// e.g: Path: "github.com/ernesto-jimenez" will only return dependents within
	// github.com/ernesto-jimenez
	Path string
}

// List looks for dependent packages tracked by godoc.org
func (g *GoDoc) List(pkg string, recursive bool) ([]string, error) {
	list := deduplist(make(map[string]bool))
	pkgs := []string{pkg}
	for len(pkgs) > 0 {
		pkg := pkgs[0]
		pkgs = pkgs[1:]
		l, err := g.fetchImportersList(pkg)
		if err != nil {
			return nil, err
		}
		list.add(l...)
		if !recursive {
			continue
		}
		l, err = g.fetchSupackages(pkg)
		if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, l...)
	}
	return list.list(), nil
}

func (g *GoDoc) url() (*url.URL, error) {
	u := g.URL
	if g.URL == "" {
		u = "https://godoc.org"
	}
	return url.Parse(u)
}

func (g *GoDoc) fetchImportersList(pkg string) ([]string, error) {
	u, err := g.url()
	if err != nil {
		return nil, err
	}
	u.Path = "/" + pkg
	u.RawQuery = "importers"
	return fetchList(u.String(), g.Path)
}

func (g *GoDoc) fetchSupackages(pkg string) ([]string, error) {
	u, err := g.url()
	if err != nil {
		return nil, err
	}
	u.Path = "/" + pkg
	return fetchList(u.String(), pkg)
}

func fetchList(url, prefix string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch %s - %s", url, res.Status)
	}
	list := []string{}
	webdevdata.ProcessMatchingTagsReader(res.Body, "table a", func(node *html.Node) {
		p := strings.TrimLeft(webdevdata.GetAttr("href", node.Attr), "/")
		if !strings.HasPrefix(p, prefix) {
			return
		}
		list = append(list, p)
	})
	return list, nil
}
