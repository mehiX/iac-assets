package gitlab

import (
	"bytes"
	"encoding/base64"

	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v3"
)

type Collector struct {
	BaseURL string
	Token   string
}

type Source struct {
	Tenant   string
	Project  string
	Ref      string
	Filepath string
}

type Result struct {
	CommitID string
	Zones    Zones
	Error    error
}

func (c *Collector) ReadFile(src Source) (*gitlab.File, error) {
	git, err := gitlab.NewClient(c.Token, gitlab.WithBaseURL(c.BaseURL))
	if err != nil {
		return nil, err
	}

	gf := &gitlab.GetFileOptions{
		Ref: gitlab.Ptr(src.Ref),
	}

	f, _, err := git.RepositoryFiles.GetFile(src.Project, src.Filepath, gf)

	return f, err
}

func (c *Collector) Query(src Source) Result {

	var zones Zones

	f, err := c.ReadFile(src)
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
