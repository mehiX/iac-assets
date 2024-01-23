package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var cmdVcloud = &cobra.Command{
	Use:  "vcloud",
	Long: "Show items managed in Virtual Cloud Directory as part of IAC",
	Run: func(cmd *cobra.Command, args []string) {
		data := vCloud.Collect(vCloudSrc...)
		if err := json.NewEncoder(os.Stdout).Encode(data); err != nil {
			log.Fatalln(err)
		}
	},
}
