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
	Use:   "test",
	Short: "Run tests for all depending packages",
	//Long: ``,
	PersistentPreRun: requireInsecure,
	Run: func(cmd *cobra.Command, args []string) {
		_, list := prepare(cmd, args)
		t, err := tester.NewTempTest()
		if verbose {
			t.Stdout = os.Stdout
			t.Stderr = os.Stderr
		}
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

func init() {
	RootCmd.AddCommand(testCmd)
	testCmd.Flags().BoolVar(&buildPkgs, "build", false, "Build packages")
}
