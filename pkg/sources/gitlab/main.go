package gitlab

import (
	"bytes"
	"encoding/base64"
	"log/slog"

	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v3"
)

type Source struct {
	BaseURL  string `json:"baseurl"`
	Token    string `json:"token"`
	Tenant   string `json:"tenant"`
	Project  string `json:"project"`
	Ref      string `json:"ref"`
	Filepath string `json:"filepath"`
}

type Result struct {
	CommitID string
	Zones    Zones
	Error    error
}

func (s Source) ReadFile() (*gitlab.File, error) {
	git, err := gitlab.NewClient(s.Token, gitlab.WithBaseURL(s.BaseURL))
	if err != nil {
		return nil, err
	}

	gf := &gitlab.GetFileOptions{
		Ref: gitlab.Ptr(s.Ref),
	}

	f, _, err := git.RepositoryFiles.GetFile(s.Project, s.Filepath, gf)

	return f, err
}

func (s Source) Query() Result {

	slog.Info("query gitlab", "tenant", s.Tenant)

	var zones Zones

	f, err := s.ReadFile()
	if err != nil {
		return Result{Error: err}
	}

	contents := f.Content
	if f.Encoding == "base64" {
		content, err := base64.StdEncoding.DecodeString(f.Content)
		if err != nil {
			return Result{Error: err}
		}
		contents = string(content)
	}

	err = yaml.NewDecoder(bytes.NewReader([]byte(contents))).Decode(&zones)

	return Result{
		CommitID: f.LastCommitID,
		Zones:    zones,
		Error:    err}
}
