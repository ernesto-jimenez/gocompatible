package cmd

import (
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gophergala2016/gocompatible/util"
	"github.com/spf13/cobra"
)

func prepare(cmd *cobra.Command, args []string) (string, []string) {
	p := "."
	if len(args) > 1 {
		log.Fatal("Too many arguments")
		log.Fatal(cmd.Short)
	} else if len(args) == 1 {
		p = args[0]
	}
	pkg, err := importPath(p)
	if err != nil {
		log.Fatal(err)
	}
	if verbose {
		log.Printf("Searching packages depending on %s", pkg)
	}
	list, err := dependents.List(pkg, recursive)
	if err != nil {
		log.Fatal(err)
	}
	return pkg, list
}

func requireInsecure(*cobra.Command, []string) {
	if godoc && !insecure {
		log.Fatal(`DANGEROUS ACTION!

Running tests from packages fetched from godoc is
dangerous. You will be basically running untrusted
code downloaded from the internet.

You should run gocompatible within a docker containr
or use the --filter option to just run tests from
trusted sources.

In order to run this command with --godoc use --godoc --insecure`)
	}
}

func importPath(p string) (string, error) {
	var clean string
	if clean = strings.TrimSuffix(p, "/..."); p != clean {
		recursive = true
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	pkg, err := build.Import(clean, cwd, build.FindOnly)
	if err == nil {
		return pkg.ImportPath, nil
	}
	if !godoc {
		return "", err
	}
	if strings.HasPrefix(clean, ".") {
		clean, _ = filepath.Abs(clean)
	}
	return util.StripGopath(clean), nil
}
