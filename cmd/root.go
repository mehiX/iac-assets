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

// defined here since they are used by more than 1 command
var gitlabDgpT *gitlab.Collector
var vCloud *vcloud.Collector

var vCloudSrc = make([]vcloud.Source, 0)

func init() {
	godotenv.Load()

	gitlabDgpT = &gitlab.Collector{
		Token:   os.Getenv("GITLAB_TOKEN"),
		BaseURL: os.Getenv("GITLAB_BASEURL"),
	}

	vCloud = &vcloud.Collector{
		Username: os.Getenv("PICARD_USER"),
		Password: os.Getenv("PICARD_PASSWORD"),
	}

	ep := strings.Split(os.Getenv("VCLOUD_ENDPOINTS"), ",")
	tenants := strings.Split(os.Getenv("VCLOUD_TENANTS"), ",")

	for _, t := range tenants {
		src := vcloud.Source{Endpoints: ep, Tenant: t}
		vCloudSrc = append(vCloudSrc, src)
	}

	cmdRoot.AddCommand(cmdGitlab, cmdVcloud, cmdServeHttp)
}

func Execute() error {
	return cmdRoot.Execute()
}
