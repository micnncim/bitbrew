package bitbrew

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/micnncim/bitbrew/github"
	"github.com/micnncim/bitbrew/plugin"
)

const (
	pluginRepo = "matryer/bitbar-plugins"
)

type Service interface {
	Search(ctx context.Context, q string) (plugin.Plugins, error)
	SearchByFilename(ctx context.Context, filename string) (plugin.Plugins, error)
	ListLocal() (plugin.Plugins, error)
	Load() error
	Save() error
	Install(plugin *plugin.Plugin) error
	Uninstall(plugin *plugin.Plugin) error
}

type service struct {
	github       github.Service
	plugins      plugin.Plugins
	formulaPath  string
	pluginFolder string
}

func NewService(gh github.Service, formulaPath, pluginFolder string) Service {
	return &service{
		github:       gh,
		formulaPath:  formulaPath,
		pluginFolder: pluginFolder,
	}
}

func (s *service) Search(ctx context.Context, q string) (plugin.Plugins, error) {
	q = fmt.Sprintf("%s bitbar title desc repo:%s", q, pluginRepo)
	return s.github.Search(ctx, q)
}

func (s *service) SearchByFilename(ctx context.Context, filename string) (plugin.Plugins, error) {
	q := fmt.Sprintf("filename:%s repo:%s", filename, pluginRepo)
	return s.github.Search(ctx, q)
}

func (s *service) ListLocal() (plugin.Plugins, error) {
	var plugins plugin.Plugins

	buf, err := ioutil.ReadFile(s.formulaPath)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(buf, &plugins); err != nil {
		return nil, err
	}

	return plugins, nil
}

func (s *service) Load() error {
	if !s.formulaExists() {
		return errors.New("formulaPath does not exist")
	}
	f, err := os.Open(s.formulaPath)
	if err != nil {
		return err
	}
	defer f.Close()
	var plugins plugin.Plugins
	if err := yaml.NewDecoder(f).Decode(&plugins); err != nil {
		return err
	}
	s.plugins = plugins
	return nil
}

func (s *service) Save() error {
	if !s.formulaExists() {
		f, err := os.Create(s.formulaPath)
		if err != nil {
			return err
		}
		defer f.Close()
		return yaml.NewEncoder(f).Encode(s.plugins)
	}
	f, err := os.Open(s.formulaPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewEncoder(f).Encode(s.plugins)
}

func (s *service) formulaExists() bool {
	_, err := os.Stat(s.formulaPath)
	return err == nil
}

func (s *service) Install(plugin *plugin.Plugin) error {
	resp, err := http.Get(plugin.GitHubRawURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(s.pluginFolder, plugin.Filename), buf, 0755); err != nil {
		return err
	}

	return nil
}

func (s *service) Uninstall(plugin *plugin.Plugin) error {
	return os.Remove(filepath.Join(s.pluginFolder, plugin.Filename))
}
