package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cmdShow = &cobra.Command{
	Use:   "show",
	Long:  "Print out IAC information gathered from defined sources",
	Short: "Print out IAC information",
	Run: func(cmd *cobra.Command, args []string) {
		if len(manager.Data) == 0 {
			fmt.Println("No data")
			return
		}

		json.NewEncoder(os.Stdout).Encode(manager.Data)
	},
}
