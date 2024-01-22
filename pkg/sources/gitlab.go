package sources

import (
	"bytes"
	"encoding/base64"

	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v3"
)

type GitlabCollector struct {
	Name     string // a name identifying this collector
	BaseURL  string
	Token    string
	Project  string
	Ref      string
	Filepath string
}

type GitlabResult struct {
	CommitID string
	Zones    Zones
	Error    error
}

func (c *GitlabCollector) ReadFile() (*gitlab.File, error) {
	git, err := gitlab.NewClient(c.Token, gitlab.WithBaseURL(c.BaseURL))
	if err != nil {
		return nil, err
	}

	gf := &gitlab.GetFileOptions{
		Ref: gitlab.Ptr(c.Ref),
	}

	f, _, err := git.RepositoryFiles.GetFile(c.Project, c.Filepath, gf)

	return f, err
}

func (c *GitlabCollector) Query() GitlabResult {

	var zones Zones

	f, err := c.ReadFile()
	if err != nil {
		return GitlabResult{Error: err}
	}

	contents := f.Content
	if f.Encoding == "base64" {
		content, err := base64.StdEncoding.DecodeString(f.Content)
		if err != nil {
			return GitlabResult{Error: err}
		}
		contents = string(content)
	}

	err = yaml.NewDecoder(bytes.NewReader([]byte(contents))).Decode(&zones)

	return GitlabResult{
		CommitID: f.LastCommitID,
		Zones:    zones,
		Error:    err}
}

type Machine struct {
	IPLastOctet  string `yaml:"ip_last_octet"`
	Tier         int    `yaml:"tier"`
	CpuCount     int    `yaml:"cpu_count"`
	CpuPerSocket int    `yaml:"cpu_per_socker"`
	MemorySizeGB int    `yaml:"memory_size_gb"`
	Disks        []Disk `yaml:"disks"`
}

type Disk struct {
	Bus    int `yaml:"bus"`
	Unit   int `yaml:"unit"`
	SizeGB int `yaml:"size_gb"`
}

type Vapp map[string]Machine

type Vapps map[string]Vapp

type Zone map[string]Vapps

type Zones map[string]Zone
