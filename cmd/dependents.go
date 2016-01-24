package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

// dependentsCmd represents the dependents command
var dependentsCmd = &cobra.Command{
	Use:   "dependents [<package>]",
	Short: "List all packages we can find depending on the given package",

	Long: `

This is very useful for:

 - Finding local packages depending on certain package.
 e.g: to consider upgrading a package.
 - Finding packages tracked by godoc.org as importing certain package.
 e.g: to help you maintain backwards compatibility of your open source package.

	`,

	Example: strings.TrimSpace(`

# Find all pkgs on GOPATH depending on .
gocompatible dependents

# Find all pkgs on GOPATH depending on ./...
gocompatible dependents ./...

# Find all pkgs on GOPATH depending on testify
gocompatible dependents github.com/stretchr/testify/...

# Find all pkgs tracked by godoc.org as importing testify
gocompatible dependents --godoc github.com/stretchr/testify/...

	`),
	Run: func(cmd *cobra.Command, args []string) {
		_, list := prepare(cmd, args)
		if verbose && len(list) == 0 {
			log.Println("None found")
		}
		for _, dep := range list {
			fmt.Println(dep)
		}
	},
}

func init() {
	RootCmd.AddCommand(dependentsCmd)
}
