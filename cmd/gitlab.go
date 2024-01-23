package cmd

import (
	"encoding/json"
	"log"
	"os"

	"git.lpc.logius.nl/logius/open/dgp/launchpad/iac-assets/pkg/sources/gitlab"
	"github.com/spf13/cobra"
)

var cmdGitlab = &cobra.Command{
	Use:  "gitlab",
	Long: "Show items managed in Gitlab as part of IAC",
	Run: func(cmd *cobra.Command, args []string) {

		src := gitlab.Source{
			Tenant:   "DGP-DGP-OT",
			Ref:      "main",
			Project:  "logius/open/dgp/infra-config-dgp-dgp-ot",
			Filepath: "infra/dgp-t/vm.yml",
		}

		if err := json.NewEncoder(os.Stdout).Encode(gitlabDgpT.Collect(src)); err != nil {
			log.Fatalln(err)
		}
	},
}
