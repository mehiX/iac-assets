package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"git.lpc.logius.nl/logius/open/dgp/launchpad/iac-assets/pkg/sources/gitlab"
	"git.lpc.logius.nl/logius/open/dgp/launchpad/iac-assets/pkg/sources/vcloud"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var configFile string

var cmdRoot = &cobra.Command{
	Use:  "iac",
	Long: "Show IAC information",
}

func init() {
	godotenv.Load()

	cmdRoot.AddCommand(cmdGitlab, cmdVcloud, cmdServeHttp)
	cmdRoot.PersistentFlags().StringVarP(&configFile, "config", "c", "config.json", "Location for the configuration file")
}

func Execute() error {
	return cmdRoot.Execute()
}

func getVCloudCollector() *vcloud.Collector {
	return &vcloud.Collector{
		Username: os.Getenv("PICARD_USER"),
		Password: os.Getenv("PICARD_PASSWORD"),
	}
}

func getGitlabCollectore() *gitlab.Collector {
	return &gitlab.Collector{
		Token:   os.Getenv("GITLAB_TOKEN"),
		BaseURL: os.Getenv("GITLAB_BASEURL"),
	}
}
func getVCloudSources() (vCloudSrc []vcloud.Source) {

	type config struct {
		VCloud map[string][]string `json:"vcloud"`
	}

	b, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var cfg config
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	ep := cfg.VCloud["endpoints"]
	tenants := cfg.VCloud["tenants"]

	for _, t := range tenants {
		src := vcloud.Source{Endpoints: ep, Tenant: t}
		vCloudSrc = append(vCloudSrc, src)
	}

	return
}

func getGitlabSources() (src []gitlab.Source) {

	type gsources struct {
		Sources []gitlab.Source `json:"sources"`
	}
	type config struct {
		Gitlab gsources `json:"gitlab"`
	}

	b, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	var cfg config
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	return cfg.Gitlab.Sources
}
