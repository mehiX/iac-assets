package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"git.lpc.logius.nl/logius/open/dgp/launchpad/iac-assets/pkg/sources/vcloud"
	"github.com/spf13/cobra"
)

var cmdVcloud = &cobra.Command{
	Use:  "vcloud",
	Long: "Show items managed in Virtual Cloud Directory as part of IAC",
	Run: func(cmd *cobra.Command, args []string) {

		src := getVCloudSources()

		if len(src) == 0 {
			fmt.Println("No sources defined. Nothing to do. Exiting")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		data := vcloud.Collect(ctx, src...)
		if err := json.NewEncoder(os.Stdout).Encode(data); err != nil {
			log.Fatalln(err)
		}
	},
}
