// Copyright © 2016 Ernesto Jimenez <me@ernesto-jimenez.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gophergala2016/gocompatible/importers"
	"github.com/spf13/cobra"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gocompatible",
	Short: "Find dependent packages and run their tests to ensure backwards compatibility",

	Long: `

Backwards compatibility is really important in Go, specially with the
limitations around package versioning.

_gocompatible_ helps you keep your packages backwards compatible:

 - Find packages that depend on yours from your GOPATH or from godoc.org
 - Run tests from all the packages depending on yours
 - Check whether any dependeng packages break on different revisions

# Security

` + securityInstructions,

	Example: strings.Join([]string{
		dependentsCmd.Example,
		testCmd.Example,
		diffCmd.Example,
	}, "\n\n"),
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	log.SetFlags(0)
	cobra.OnInitialize(initDependents)

	RootCmd.PersistentFlags().BoolVarP(&local, "local", "l", true, "find dependents in your $GOPATH")
	RootCmd.PersistentFlags().BoolVar(&godoc, "godoc", false, "find dependents in godoc.org")
	RootCmd.PersistentFlags().StringVarP(&filter, "filter", "f", "", "select dependents within the given path")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().BoolVarP(&recursive, "recurisve", "r", false, "check subpackages too")
}

var (
	local      bool
	godoc      bool
	verbose    bool
	filter     string
	recursive  bool
	insecure   bool
	dependents importers.Lister
)

func initDependents() {
	if godoc {
		dependents = &importers.GoDoc{Path: filter}
	} else {
		dependents = &importers.Local{Path: filter}
	}
}

var securityInstructions = strings.TrimSpace(`

Never run _test_ or _diff_ commands with --godoc
in your machine. This will download untrusted code
and run it, which is very dangerous.

You must always run those commands within an
issolated container. We've published a docker
image you can use to do it quickly.

Example:

docker run --rm quay.io/ernesto_jimenez/gocompatible \
  gocompatible test github.com/stretchr/testify/assert \
  --filter github.com/uber \
  --godoc --insecure

`)
