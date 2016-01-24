package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// dependentsCmd represents the dependents command
var dependentsCmd = &cobra.Command{
	Use:   "dependents [<package>]",
	Short: "List all packages we can find depending on the given package",
	//Long:  ``,
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
