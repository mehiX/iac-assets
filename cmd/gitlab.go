package cmd

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"git.lpc.logius.nl/logius/open/dgp/launchpad/iac-assets/pkg/sources/gitlab"
	"github.com/spf13/cobra"
)

var cmdGitlab = &cobra.Command{
	Use:  "gitlab",
	Long: "Show items managed in Gitlab as part of IAC",
	Run: func(cmd *cobra.Command, args []string) {

		src := getGitlabSources()

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := json.NewEncoder(os.Stdout).Encode(gitlab.Collect(ctx, src...)); err != nil {
			log.Fatalln(err)
		}
	},
}
