package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

const (
	defaultRootConfig = ".gimme/config.yaml"
)

type Config struct {
	Dryrun         bool   `yaml:",omitempty"`
	Manifest       bool   `yaml:",omitempty"`
	HomeOverride   string `yaml:"homeOverride"`
	KubeconfigEdit bool   `yaml:"kubeconfigEdit"`
}

func LoadRootConfig(dryrun, manifest bool) (Config, error) {
	configPath := path.Join(os.Getenv("HOME"), defaultRootConfig)
	override, set := os.LookupEnv("GIMME_CONFIG_PATH")
	if set {
		configPath = override
	}

	base := Config{
		Dryrun:   dryrun,
		Manifest: manifest,
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return base, fmt.Errorf("failed to read config file: %s", err)
	}

	err = yaml.Unmarshal(data, &base)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal config yaml: %s", err)
	}

	return base, err
}

func (c *Config) DetermineHome() string {
	home := os.Getenv("HOME")
	if c.HomeOverride != "" {
		home = c.HomeOverride
	}
	return home
}

func (c *Config) GetConfigFilePath() string {
	return path.Join(c.DetermineHome(), ".gimme", "config.yaml")
}

func (c *Config) GetAliasFilePath() string {
	return path.Join(c.DetermineHome(), ".gimme", "aliases.yaml")
}

func (c *Config) EnsureSetup() error {
	home := os.Getenv("HOME")
	if len(c.HomeOverride) > 0 {
		home = c.HomeOverride
	}
	gimmeConfDir := path.Join(home, ".gimme")
	err := os.Mkdir(gimmeConfDir, 775)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		return fmt.Errorf("failed to create .gimme config directory: %s", err)
	}

	gimmeConfFile := path.Join(gimmeConfDir, "config.yaml")

	if _, err := os.Stat(gimmeConfFile); errors.Is(err, fs.ErrNotExist) {
		emptyData := make([]byte, 0)

		err := os.WriteFile(gimmeConfFile, emptyData, 775)
		if err != nil {
			return fmt.Errorf("failed to create empty config file: %s", err)
		}
	}

	return nil
}
