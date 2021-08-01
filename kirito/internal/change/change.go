package change

import (
	"fmt"
	"os"

	"github.com/go-kirito/cmd/kirito/internal/base"
	"github.com/spf13/cobra"
)

// CmdChange is kratos change log tool
var CmdChange = &cobra.Command{
	Use:   "changelog",
	Short: "Get a kirito change log",
	Long:  "Get a kirito release or commits info. Example: kirito changelog dev or kirito changelog {version}",
	Run:   run,
}
var token string
var repoURL string

func init() {
	if repoURL = os.Getenv("KIRITO_REPO"); repoURL == "" {
		repoURL = "https://github.com/go-kirito/cmd.git"
	}
	CmdChange.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "github repo")
	token = os.Getenv("GITHUB_TOKEN")
}

func run(cmd *cobra.Command, args []string) {
	owner, repo := base.ParseGithubUrl(repoURL)
	api := base.GithubApi{Owner: owner, Repo: repo, Token: token}
	version := "latest"
	if len(args) > 0 {
		version = args[0]
	}
	if version == "dev" {
		info := api.GetCommitsInfo()
		fmt.Print(base.ParseCommitsInfo(info))
	} else {
		info := api.GetReleaseInfo(version)
		fmt.Print(base.ParseReleaseInfo(info))
	}
}
