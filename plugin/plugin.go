package plugin

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
