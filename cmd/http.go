package cmd

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

const defaultAddr = "127.0.0.1:8080"

var cmdServeHttp = &cobra.Command{
	Use:   "serve",
	Long:  "Serve results over HTTP",
	Short: "Server results over HTTP",
	Run: func(cmd *cobra.Command, args []string) {
		addr := defaultAddr
		if len(args) > 0 {
			addr = args[0]
		}

		srvr := http.Server{
			Addr:    addr,
			Handler: handler(),
		}

		fmt.Printf("Listening on %s\n", srvr.Addr)
		if err := srvr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err.Error())
		}
	},
}

//go:embed tmpl
var htmlTemplates embed.FS

func handler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/json", handleGetJsonData)
	m.HandleFunc("/html", handleGetHtmlData())
	return m
}

func handleGetJsonData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(manager.Data); err != nil {
		log.Printf("Encoding Json response: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleGetHtmlData() http.HandlerFunc {
	tmpl, err := template.ParseFS(htmlTemplates, "**/*.tmpl")
	if err != nil {
		log.Fatalln("Cannot load html templates")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")

		if err := tmpl.ExecuteTemplate(w, "main", manager.Data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
