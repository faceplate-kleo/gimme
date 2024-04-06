package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func DiscoverPaths() ([]string, error) {
	home := os.Getenv("HOME")
	paths := make([]string, 0)

	regex, err := regexp.Compile(`\.gimme\.yaml$`)
	if err != nil {
		return paths, fmt.Errorf("failed to compile regex: %s", err)
	}

	err = filepath.Walk(home, func(path string, info os.FileInfo, errInternal error) error {
		if errInternal == nil && regex.MatchString(info.Name()) {
			println(path)
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return paths, fmt.Errorf("failed to walk directory hierarchy: %s", err)
	}

	return paths, nil
}
