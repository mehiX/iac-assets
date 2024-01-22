package cmd

import (
	"iac-gitlab-assets/pkg/iac"
	"iac-gitlab-assets/pkg/sources"
	"os"

	"github.com/spf13/cobra"
)

var manager = iac.NewManager()

var cmdRoot = &cobra.Command{
	Use:  "iac",
	Long: "Show items managed in Gitlab as part of IAC",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		manager.Collect(gitlabDgpT)
	},
}

var gitlabDgpT = &sources.GitlabCollector{
	Name:     "Gitlab DGP-T",
	Token:    os.Getenv("GITLAB_TOKEN"),
	BaseURL:  "https://git.lpc.logius.nl",
	Ref:      "main",
	Project:  "logius/open/dgp/infra-config-dgp-dgp-ot",
	Filepath: "infra/dgp-t/vm.yml",
}

func init() {
	cmdRoot.AddCommand(cmdShow, cmdServeHttp)
}

func Execute() error {
	return cmdRoot.Execute()
}
