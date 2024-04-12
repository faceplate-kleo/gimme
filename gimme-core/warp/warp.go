package warp

import (
	"fmt"
	"github.com/faceplate-kleo/gimme-core/config"
	"github.com/faceplate-kleo/gimme-core/discovery"
	"path"
)

func Warp(alias string, conf *config.Config) (string, error) {
	if conf == nil {
		return "", fmt.Errorf("failed to warp: config reference is nil")
	}

	aliases, err := discovery.ReadAliasFile(conf)
	if err != nil {
		return "", fmt.Errorf("failed to read alias file: %s", err)
	}

	gimmePath, ok := aliases[alias]
	if !ok {
		return "", fmt.Errorf("failed to warp: alias %s not found", alias)
	}

	if conf.Manifest {
		fmt.Printf("[WARP] %s\n", path.Dir(gimmePath))
	} else {
		fmt.Printf("If I had the power, I would have warped you to %s\n", path.Dir(gimmePath))
	}

	return gimmePath, nil
}
