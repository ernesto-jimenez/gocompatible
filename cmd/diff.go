package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gophergala2016/gocompatible/tester"
	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use: `diff --from <commit> [<package>] [flags]
  gocompatible diff <pull-request-url>`,
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
		args, remotes, err := fillPRDetails(args)
		if err != nil {
			log.Fatal(err)
		}
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
		if len(remotes) == 2 {
			debugf("Fetch remote %s", remotes[0])
			t.FetchRemote(pkg, "prhead", remotes[0])
			debugf("Fetch remote %s", remotes[1])
			t.FetchRemote(pkg, "prbase", remotes[1])
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

var matchPR = regexp.MustCompile("https?://(?:www.)?(github.com/([^/]+/[^/]+))/pull/([0-9]+)")

func fillPRDetails(args []string) ([]string, []string, error) {
	if len(args) != 1 {
		return args, nil, nil
	}
	matches := matchPR.FindStringSubmatch(args[0])
	if len(matches) < 4 {
		return args, nil, nil
	}
	pkg, repo, pr := matches[1], matches[2], matches[3]
	debugf("Fetching PR info from %s", pkg)
	prURL := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%s", repo, pr)
	debugf("GET %s", prURL)
	res, err := http.Get(prURL)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	info := struct {
		Message string `json:"message"`
		Head    struct {
			Sha  string `json:"sha"`
			Repo struct {
				CloneURL string `json:"clone_url"`
			} `json:"repo"`
		} `json:"head"`
		Base struct {
			Sha  string `json:"sha"`
			Repo struct {
				CloneURL string `json:"clone_url"`
			} `json:"repo"`
		} `json:"base"`
	}{}
	if info.Message != "" {
		return nil, nil, fmt.Errorf("API Message: %s", info.Message)
	}
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return nil, nil, err
	}
	from = info.Base.Sha
	debugf("From: %s", from)
	to = info.Head.Sha
	debugf("To: %s", to)
	recursive = true
	return []string{pkg}, []string{info.Head.Repo.CloneURL, info.Base.Repo.CloneURL}, nil
}

func init() {
	RootCmd.AddCommand(diffCmd)
	diffCmd.Flags().BoolVar(&insecure, "insecure", false, "allows running testing packages from godoc")
	diffCmd.Flags().StringVarP(&from, "from", "c", "", "commit/tag/branch of <package> to checkout and select only the packages with green builds")
	diffCmd.Flags().StringVarP(&to, "to", "t", "master", "commit/tag/branch of the tested package we want to compare with from to see how many packages broke")
}
