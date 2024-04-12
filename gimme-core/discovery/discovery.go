package discovery

import (
	"fmt"
	"github.com/faceplate-kleo/gimme-core/config"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
)

func DiscoverPaths(rootPath string) ([]string, error) {
	paths := make([]string, 0)

	regex, err := regexp.Compile(`\.gimme\.yaml$`)
	if err != nil {
		return paths, fmt.Errorf("failed to compile regex: %s", err)
	}

	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, errInternal error) error {
		if errInternal == nil && regex.MatchString(info.Name()) {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return paths, fmt.Errorf("failed to walk directory hierarchy: %s", err)
	}

	return paths, nil
}

func PopulateAliases(paths []string) (map[string]string, error) {
	aliases := make(map[string]string)
	conf := config.Config{
		Dryrun:   false,
		Manifest: false,
	}

	for _, path := range paths {
		conf, err := config.ReadGimmeFileV1(path, &conf)
		if err != nil {
			return aliases, fmt.Errorf("failed to load config %s for alias information: %s", path, err)
		}
		aliases[conf.Gimme.Alias] = path
	}

	return aliases, nil
}

func SyncAliasFile(conf *config.Config) error {
	if conf == nil {
		return fmt.Errorf("failed to sync alias file: config reference is nil")
	}

	home := conf.DetermineHome()
	paths, err := DiscoverPaths(home)
	if err != nil {
		return fmt.Errorf("failed to discover gimme paths: %s", err)
	}
	aliases, err := PopulateAliases(paths)
	if err != nil {
		return fmt.Errorf("failed to populate aliases: %s", err)
	}

	yamlData, err := yaml.Marshal(aliases)
	if err != nil {
		return fmt.Errorf("failed to marshal alias yaml: %s", err)
	}

	err = os.WriteFile(conf.GetAliasFilePath(), yamlData, 0664)
	if err != nil {
		return fmt.Errorf("failed to write alias yaml: %s", err)
	}

	return nil
}

func ReadAliasFile(conf *config.Config) (map[string]string, error) {
	aliases := make(map[string]string)
	if conf == nil {
		return aliases, fmt.Errorf("failed to read alias file: config reference is nil")
	}

	data, err := os.ReadFile(conf.GetAliasFilePath())
	if err != nil {
		return aliases, fmt.Errorf("failed to read alias file: %s", err)
	}

	err = yaml.Unmarshal(data, &aliases)
	if err != nil {
		return aliases, fmt.Errorf("failed to unmarshal alias yaml: %s", err)
	}

	return aliases, nil

}

func ListAliases(conf *config.Config) error {
	aliases, err := ReadAliasFile(conf)
	if err != nil {
		return fmt.Errorf("failed to list aliases: %s", err)
	}

	for alias, dest := range aliases {
		fmt.Printf("%-20s\t%s\n", alias, dest)
	}
	return nil
}
