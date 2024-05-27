package discovery

import (
	"fmt"
	"github.com/faceplate-kleo/gimme-core/config"
	"gopkg.in/yaml.v3"
	"os"
	"path"
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

func WriteAliasFile(conf *config.Config, aliases map[string]string) error {
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

	return WriteAliasFile(conf, aliases)
}

func AppendToAliasFile(conf *config.Config, toAppend map[string]string) error {
	if conf == nil {
		return fmt.Errorf("failed to append to alias file: config reference is nil")
	}

	aliases, err := ReadAliasFile(conf)
	if err != nil {
		return fmt.Errorf("failed to read alias file: %s", err)
	}

	for k, v := range toAppend {
		aliases[k] = v
	}

	return WriteAliasFile(conf, aliases)
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
		fmt.Printf("%-20s\t%s\n", alias, path.Dir(dest))
	}
	return nil
}

func ValidateAlias(conf *config.Config, alias string) (bool, error) {
	aliases, err := ReadAliasFile(conf)
	if err != nil {
		return false, fmt.Errorf("failed to list aliases: %s", err)
	}

	for _, check := range aliases {
		if check == alias {
			return false, nil
		}
	}
	return true, nil
}

func PinDirectory(conf *config.Config, path, alias string) error {
	if ok, err := ValidateAlias(conf, alias); !ok {
		if err != nil {
			return fmt.Errorf("failed to validate alias %q: %s", alias, err)
		}
		return fmt.Errorf("alias already exists: %s", alias)
	}

	blankConfig := config.NewSpecV1(alias, nil)

	yamlData, err := yaml.Marshal(blankConfig)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write gimmefilefile %s: %s", path, err)
	}

	return AppendToAliasFile(conf, map[string]string{alias: path})
}
