package config

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Config represent a config file for CLI
type Config struct {
	*BitBar
	*GitHub
}

// BitBar represent a BitBar config for CLI
type BitBar struct {
	FormulaPath  string `mapstructure:"formula_path"`
	PluginFolder string `mapstructure:"plugin_folder"`
}

// GitHub represent a GitHub config for CLI
type GitHub struct {
	Token string
}

var (
	newDefaultViperFunc = newDefaultViper
	initConfigFunc      = initConfig
	configDirFunc       = configDir
)

var (
	defaultConfigDir  = ".config/bitbrew"
	defaultConfigName = "config.yaml"
)

const (
	DefaultGitHubToken = "<GITHUB_ACCESS_TOKEN>"
)

// New instantiate Config
func New() (*Config, error) {
	configDir, err := configDirFunc()
	if err != nil {
		return nil, err
	}
	conf, err := initConfigFunc(configDir)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(conf.BitBar.PluginFolder, 0755); err != nil {
		return nil, err
	}
	return conf, nil
}

func newDefaultViper(configDir string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath(configDir)

	v.SetDefault("bitbar.formula_path", filepath.Join(configDir, "formula.yaml"))
	v.SetDefault("bitbar.plugin_folder", filepath.Join(configDir, "plugins"))
	v.SetDefault("github.token", DefaultGitHubToken)

	return v, nil
}

func initConfig(configDir string) (*Config, error) {
	v, err := newDefaultViperFunc(configDir)
	if err != nil {
		return nil, err
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := os.MkdirAll(configDir, 0755); err != nil {
				return nil, err
			}
			if err := v.WriteConfigAs(filepath.Join(configDir, defaultConfigName)); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, defaultConfigDir), nil
}

// Edit open a config file with editor
func Edit() error {
	editor := getEditor()
	if editor == "" {
		return errors.New("require $EDITOR or vim")
	}
	configDir, err := configDir()
	if err != nil {
		return err
	}
	return runEditor(editor, filepath.Join(configDir, defaultConfigName))
}

func getEditor() string {
	if env := os.Getenv("EDITOR"); env != "" {
		return env
	}
	p, err := exec.LookPath("vim")
	if err != nil {
		return ""
	}
	return p
}

func runEditor(editor, path string) error {
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to execute %s", editor)
	}
	return nil
}
