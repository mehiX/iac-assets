package cmd

import (
	"os"
	"strings"

	"git.lpc.logius.nl/logius/open/dgp/launchpad/iac-assets/pkg/sources/gitlab"
	"git.lpc.logius.nl/logius/open/dgp/launchpad/iac-assets/pkg/sources/vcloud"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var cmdRoot = &cobra.Command{
	Use:  "iac",
	Long: "Show IAC information",
}

var gitlabDgpT *gitlab.Collector
var vCloud *vcloud.Collector

func init() {
	godotenv.Load()

	gitlabDgpT = &gitlab.Collector{
		Name:     "Gitlab DGP-T",
		Token:    os.Getenv("GITLAB_TOKEN"),
		BaseURL:  "https://git.lpc.logius.nl",
		Ref:      "main",
		Project:  "logius/open/dgp/infra-config-dgp-dgp-ot",
		Filepath: "infra/dgp-t/vm.yml",
	}

	vCloud = &vcloud.Collector{
		Username:  os.Getenv("PICARD_USER"),
		Password:  os.Getenv("PICARD_PASSWORD"),
		Endpoints: strings.Split(os.Getenv("VCLOUD_ENDPOINTS"), ","),
		Tenants:   strings.Split(os.Getenv("VCLOUD_TENANTS"), ","),
	}

	cmdRoot.AddCommand(cmdGitlab, cmdVcloud, cmdServeHttp)
}

func Execute() error {
	return cmdRoot.Execute()
}
