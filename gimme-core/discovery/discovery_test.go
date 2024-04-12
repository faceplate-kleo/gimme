package discovery

import (
	"github.com/faceplate-kleo/gimme-core/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestWalk(t *testing.T) {
	_, err := DiscoverPaths(os.Getenv("HOME"))
	assert.Nil(t, err)
}

func TestSyncAliases(t *testing.T) {
	err := os.Setenv("GIMME_CONFIG_PATH", "../test/.gimme/test-config.yaml")
	assert.Nil(t, err, "error from os.Setenv: %s", err)

	conf, err := config.LoadRootConfig(false, false)
	assert.Nil(t, err, "error from config.LoadRootConfig: %s", err)

	err = SyncAliasFile(&conf)
	assert.Nil(t, err, "error from discovery.SyncAliasFile %s", err)
}
