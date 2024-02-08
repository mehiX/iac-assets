package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"git.lpc.logius.nl/logius/open/dgp/launchpad/iac-assets/pkg/sources/gitlab"
	"git.lpc.logius.nl/logius/open/dgp/launchpad/iac-assets/pkg/sources/vcloud"

	"github.com/spf13/cobra"
)

var config Config

var cmdRoot = &cobra.Command{
	Use:  "iac",
	Long: "Show IAC information",
}

func init() {
	cmdRoot.AddCommand(cmdGitlab, cmdVcloud, cmdServeHttp)
	configFile := cmdRoot.PersistentFlags().StringP("config", "c", "config.json", "Location for the configuration file")

	b, err := os.ReadFile(filepath.Clean(*configFile))
	if err != nil {
		log.Println("Reading initial config: ", err)
	} else {
		if err := json.NewDecoder(bytes.NewReader(b)).Decode(&config); err != nil {
			log.Println("Decoding initial config: ", err)
		}
	}
}

func Execute() error {
	return cmdRoot.Execute()
}

func getVCloudSources() []vcloud.Source {

	return config.VCloud.Sources
}

func getGitlabSources() []gitlab.Source {

	return config.Gitlab.Sources
}

type Config struct {
	Gitlab GlabSources `json:"gitlab"`
	VCloud VcSources   `json:"vcloud"`
}

func (c Config) IsValid() bool {
	return len(c.Gitlab.Sources) != 0 || len(c.VCloud.Sources) != 0
}

type VcSources struct {
	Sources []vcloud.Source `json:"sources"`
}

type GlabSources struct {
	Sources []gitlab.Source `json:"sources"`
}
