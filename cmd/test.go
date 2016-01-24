package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gophergala2016/gocompatible/tester"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test [options] [<commit>] [<package>]",
	Short: "Run tests for all depending packages",
	Long: `

This can be useful to consider whether to upgrade one of your dependencies
across your organisation.

e.g.: Imagine you work at uber and whant to check your packages running testify
will work with the latests version on master. You could run:

gocompatible githb.com/stretchr/testify/... \
  --filter github.com/uber

This will find all the packages on your GOPATH depending on testify and run
their tests with the latest version of testify.

You could also checkout a specific branch or tag from testify:

gocompatible githb.com/stretchr/testify/... \
  --checkout v1.1
  --filter github.com/uber

	`,

	Example: strings.TrimSpace(`

# Run test for all packages in GOPATH dependent on testify/assert
gocompatible test github.com/stretchr/testify/assert

# Run tests from all uber packages tracked by godoc.org as dependent on testify/assert
gocompatible test github.com/stretchr/testify/assert \
  --filter github.com/uber \
  --godoc --insecure

	`),
	PersistentPreRun: requireInsecure,
	Run: func(cmd *cobra.Command, args []string) {
		pkg, list := prepare(cmd, args)
		t, err := tester.NewTempTest()
		if err != nil {
			log.Fatal(err)
		}
		if verbose {
			t.Stdout = os.Stdout
			t.Stderr = os.Stderr
		}
		err = t.Get(pkg)
		if err != nil {
			log.Fatal(err)
		}
		err = t.Checkout(pkg, checkout)
		if err != nil {
			log.Fatal(err)
		}
		res := tester.Run(t, list, os.Stdout)
		if err := res.Err(); err != nil {
			log.Fatal(err)
		}
		fmt.Println(res.String())
	},
}

var buildPkgs bool

var checkout string

func init() {
	RootCmd.AddCommand(testCmd)
	testCmd.Flags().BoolVar(&buildPkgs, "build", false, "Build packages")
	testCmd.Flags().BoolVar(&insecure, "insecure", false, "Allows running testing packages from godoc")
	testCmd.Flags().StringVarP(&checkout, "checkout", "c", "master", "The commit/tag/branch of the tested package we want to checkout")
}
