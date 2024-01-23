package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var cmdGitlab = &cobra.Command{
	Use:  "gitlab",
	Long: "Show items managed in Gitlab as part of IAC",
	Run: func(cmd *cobra.Command, args []string) {

		coll := getGitlabCollectore()
		src := getGitlabSources()

		if err := json.NewEncoder(os.Stdout).Encode(coll.Collect(src...)); err != nil {
			log.Fatalln(err)
		}
	},
}
