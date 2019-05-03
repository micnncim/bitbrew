package plugin

import (
	"gopkg.in/yaml.v2"
)

type Plugin struct {
	Name         string `yaml:",omitempty"`
	Filename     string `yaml:",omitempty"`
	Description  string `yaml:"-"`
	Path         string `yaml:",omitempty"`
	BitBarURL    string `yaml:",omitempty"`
	GitHubURL    string `yaml:",omitempty"`
	GitHubRawURL string `yaml:",omitempty"`
}

type Plugins []*Plugin

func (p *Plugins) Marshal() ([]byte, error) {
	return yaml.Marshal(p)
}
