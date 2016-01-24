package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gophergala2016/gocompatible/tester"
	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff [options] --from <commit> [<package>]",
	Short: "Compare what depending packages break between two diffent revisions",
	Long: `

This is incredibly useful to find whether your new changes will break any
dependents.

_gocompatible_ will checkout the --from version, run the tests, and run the
tests again for --to.

Packages failing in --from are discarded from the --to run, so _gocompatible_
will only fail if the changes in --to break any package that worked in --from.

	`,
	Example: strings.TrimSpace(`

# Check whether changes between testify v1.0 and
# v1.1.1 broke any tests in aws-sdk-go
gocompatible diff github.com/stretchr/testify/... \
  --filter github.com/aws/aws-sdk-go/awstesting \
  --from v1.0 --to v1.1.1

	`),
	PreRun: requireInsecure,
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
		log.Printf("%d packages skipped due to tests failing in %s", skipped, from)
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
	diffCmd.Flags().BoolVar(&insecure, "insecure", false, "allows running testing packages from godoc")
	diffCmd.Flags().StringVarP(&from, "from", "c", "", "commit/tag/branch of <package> to checkout and select only the packages with green builds")
	diffCmd.Flags().StringVarP(&to, "to", "t", "master", "commit/tag/branch of the tested package we want to compare with from to see how many packages broke")
}
