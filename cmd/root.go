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

	"github.com/gophergala2016/gocompatible/importers"
	"github.com/spf13/cobra"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "gocompatible",
	//Short: "A brief description of your application",
	//Long: `A longer description that spans multiple lines and likely contains
	//examples and usage of using your application. For example:

	//Cobra is a CLI library for Go that empowers applications.
	//This application is a tool to generate the needed files
	//to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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

	RootCmd.PersistentFlags().BoolVarP(&local, "local", "l", true, "Find dependents in your $GOPATH")
	RootCmd.PersistentFlags().BoolVar(&godoc, "godoc", false, "Find dependents in godoc.org")
	RootCmd.PersistentFlags().StringVarP(&filter, "filter", "f", "", "Filter dependents within the given path")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	RootCmd.PersistentFlags().BoolVarP(&recursive, "recurisve", "r", false, "Check subpackages too")
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
