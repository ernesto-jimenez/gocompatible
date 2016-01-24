package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/gophergala2016/gocompatible/tester"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test [options] [<commit>] [<package>]",
	Short: "Run tests for all depending packages",
	//Long: ``,
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
