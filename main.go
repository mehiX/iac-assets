package main

import (
	"log"

	"git.lpc.logius.nl/logius/open/dgp/launchpad/iac-assets/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
