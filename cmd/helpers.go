package cmd

import (
	"log"

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
