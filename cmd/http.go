package cmd

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"git.lpc.logius.nl/logius/open/dgp/launchpad/iac-assets/pkg/sources/gitlab"
	"git.lpc.logius.nl/logius/open/dgp/launchpad/iac-assets/pkg/sources/vcloud"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
)

const defaultAddr = "127.0.0.1:8080"

var cmdServeHttp = &cobra.Command{
	Use:   "serve",
	Long:  "Serve results over HTTP",
	Short: "Serve results over HTTP",
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
	m := chi.NewRouter()
	m.Get("/", handleHome)
	m.Get("/{src}/json", handleGetJsonData)
	m.Get("/{src}/html", handleGetHtmlData())
	return m
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/gitlab/html", http.StatusFound)
}

func handleGetJsonData(w http.ResponseWriter, r *http.Request) {

	src := chi.URLParam(r, "src")
	var data any
	switch src {
	case "gitlab":
		src := getGitlabSources()
		data = gitlab.Collect(src...)
	case "vcloud":
		src := getVCloudSources()
		data = vcloud.Collect(src...)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
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

		src := chi.URLParam(r, "src")
		var data any
		switch src {
		case "gitlab":
			src := getGitlabSources()
			data = gitlab.Collect(src...)
		case "vcloud":
			src := getVCloudSources()
			data = vcloud.Collect(src...)
		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := tmpl.ExecuteTemplate(w, src+"_main", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
