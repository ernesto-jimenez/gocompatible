package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"

	"github.com/gophergala2016/gocompatible/importers"
	"github.com/gophergala2016/gocompatible/tester"
)

var (
	pkg     = flag.String("pkg", ".", "what package to check dependents")
	subpkgs = flag.Bool("subpkgs", false, "wether to search dependents from subpackages too")
	godoc   = flag.Bool("godoc", false, "fetch dependents from godoc.org instead of locally")
	inPath  = flag.String("in-path", "", "scope dependencies within the given path")
	print   = flag.Bool("print", false, "just print the list of dependents")
	verbose = flag.Bool("v", false, "verbose output")

	i importers.Lister
)

func init() {
	flag.Parse()
	if *godoc {
		i = &importers.GoDoc{Path: *inPath}
	} else {
		i = &importers.Local{Path: *inPath}
	}
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	pkg, err := build.Import(*pkg, cwd, build.FindOnly)
	if err != nil {
		log.Fatal(err)
	}
	if *verbose {
		fmt.Printf("Running gocompatible for %s\n", pkg.ImportPath)
	}
	list, err := i.List(pkg.ImportPath, *subpkgs)
	if err != nil {
		log.Fatal(err)
	}
	if *print {
		for _, l := range list {
			fmt.Println(l)
		}
		os.Exit(0)
	}
	t, err := tester.NewTempTest()
	if *verbose {
		t.Stdout = os.Stdout
		t.Stderr = os.Stderr
	}
	if err != nil {
		log.Fatal(err)
	}
	for _, pkg := range list {
		if err := t.Get(pkg); err != nil {
			fmt.Printf("%-7s %s - %s: %s\n", "FAIL", pkg, "go get", err)
			continue
		}
		if err := t.Build(pkg); err != nil {
			fmt.Printf("%-7s %s - %s: %s\n", "FAIL", pkg, "go build", err)
			continue
		}
		if err := t.Test(pkg); err != nil {
			if !*verbose {
				fmt.Printf("%-7s %s - %s: %s\n", "FAIL", pkg, "go test", err)
			}
			continue
		}
		if !*verbose {
			fmt.Printf("%-7s %s\n", "ok", pkg)
		}
	}
}
