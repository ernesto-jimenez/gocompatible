package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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

	PersistentPreRun: func(command *cobra.Command, _ []string) {
		if docker {
			cmd := exec.Command("docker", strings.Split(
				"run --rm quay.io/ernesto_jimenez/gocompatible:latest gocompatible", " ",
			)...)
			for _, arg := range os.Args[1:] {
				if arg != "--docker" && arg != "-d" {
					cmd.Args = append(cmd.Args, arg)
				}
			}
			if command == testCmd || command == diffCmd {
				cmd.Args = append(cmd.Args, "--insecure")
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}
	},
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
	RootCmd.PersistentFlags().BoolVarP(&docker, "docker", "d", false, "run the command in a docker container")
}

var (
	local      bool
	godoc      bool
	verbose    bool
	filter     string
	recursive  bool
	insecure   bool
	docker     bool
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

gocompatible test github.com/stretchr/testify/assert \
  --filter github.com/uber \
  --godoc --docker

Which is equivalent to:

docker run --rm quay.io/ernesto_jimenez/gocompatible:latest \
  gocompatible test github.com/stretchr/testify/assert \
  --filter github.com/uber \
  --godoc --insecure

`)
