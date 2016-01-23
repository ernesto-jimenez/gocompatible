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
	pkg       = flag.String("pkg", ".", "what package to check dependents")
	subpkgs   = flag.Bool("subpkgs", false, "wether to search dependents from subpackages too")
	godoc     = flag.Bool("godoc", false, "fetch dependents from godoc.org instead of locally")
	localPath = flag.String("local-path", "", "whether to scope local search to certain paths")
	print     = flag.Bool("print", false, "just print the list of dependents")

	i importers.Lister
)

func init() {
	flag.Parse()
	if *godoc {
		i = &importers.GoDoc{}
	} else {
		i = &importers.Local{Path: *localPath}
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
	fmt.Printf("Importers for %s\n", pkg.ImportPath)
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
	if err != nil {
		log.Fatal(err)
	}
	for _, pkg := range list {
		fmt.Println("##", pkg)
		fmt.Printf(" - Getting... ")
		if err := t.Get(pkg); err != nil {
			fmt.Println("!!!ERROR", err)
			continue
		}
		fmt.Println("OK")
		fmt.Printf(" - Building... ")
		if err := t.Build(pkg); err != nil {
			fmt.Println("!!!ERROR", err)
			continue
		}
		fmt.Println("OK")
		fmt.Printf(" - Testing... ")
		if err := t.Test(pkg); err != nil {
			fmt.Println("!!!ERROR", err)
			continue
		}
		fmt.Println("OK")
	}
}
