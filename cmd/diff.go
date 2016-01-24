package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/gophergala2016/gocompatible/tester"
	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:              "diff [options] --from <commit> [<package>]",
	Short:            "Compare what depending packages break between two diffent revisions",
	Long:             ``,
	PersistentPreRun: requireInsecure,
	Run: func(cmd *cobra.Command, args []string) {
		if from == "" {
			log.Println("You must specify a --from")
			log.Fatal(cmd.UsageString())
		}
		pkg, list := prepare(cmd, args)
		t, err := tester.NewTempTest()
		if err != nil {
			log.Fatal(err)
		}
		if verbose {
			t.Stdout = os.Stdout
			t.Stderr = os.Stderr
		}
		debugln("Installing package")
		err = t.Get(pkg)
		if err != nil {
			log.Fatal(err)
		}
		debugf("Checkout %s", from)
		err = t.Checkout(pkg, from)
		if err != nil {
			log.Fatal(err)
		}

		t.Stdout = ioutil.Discard
		t.Stderr = ioutil.Discard
		debugf("Runin tests from %s", from)
		res := tester.Run(t, list, ioutil.Discard)

		if verbose {
			t.Stdout = os.Stdout
			t.Stderr = os.Stderr
		}
		skipped := res.TotalFailed()

		debugf("Checkout %s", to)
		err = t.Checkout(pkg, to)
		if err != nil {
			log.Fatal(err)
		}
		debugf("Runin tests from %s", to)
		res = tester.Run(t, res.OK, os.Stdout)
		log.Printf("%d skipped due to tests failing in %s", skipped, from)
		if err := res.Err(); err != nil {
			log.Fatal(err)
		}
		log.Println(res.String())
	},
}

var (
	from string
	to   string
)

func init() {
	RootCmd.AddCommand(diffCmd)
	diffCmd.Flags().BoolVar(&buildPkgs, "build", false, "Build packages")
	diffCmd.Flags().BoolVar(&insecure, "insecure", false, "Allows running testing packages from godoc")
	diffCmd.Flags().StringVarP(&from, "from", "c", "", "The commit/tag/branch of <package> to checkout and select only the packages with green builds")
	diffCmd.Flags().StringVarP(&to, "to", "t", "master", "The commit/tag/branch of the tested package we want to compare with from to see how many packages broke")
}
