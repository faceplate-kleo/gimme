package main

import (
	"os"
	"testing"

	"github.com/faceplate-kleo/gimme/config"
	"github.com/faceplate-kleo/gimme/discovery"
	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	data, err := config.ReadGimmeFileV1(".gimme.yaml")
	if err != nil {
		t.Fatalf("Fatal: %s\n", err)
	}

	err = data.Process()
	if err != nil {
		t.Fatalf("Fatal: %s\n", err)
	}
}

func TestWalk(t *testing.T) {
	_, err := discovery.DiscoverPaths(os.Getenv("HOME"))
	assert.Nil(t, err)
}

func TestLoadConfig(t *testing.T) {
	conf, err := config.LoadRootConfig()
	assert.NotNil(t, err)
	assert.Empty(t, conf)

	err = os.Setenv("GIMME_CONFIG_PATH", "examples/example_root_config.yaml")
	assert.Nil(t, err)
	conf, err = config.LoadRootConfig()
	assert.Nil(t, err)

	assert.Equal(t, "/home/you/other/home", conf.HomeOverride)
	assert.Equal(t, true, conf.KubeconfigEdit)
}

func TestSyncAliases(t *testing.T) {
	err := os.Setenv("GIMME_CONFIG_PATH", "test/.gimme/test-config.yaml")
	assert.Nil(t, err)

	conf, err := config.LoadRootConfig()
	assert.Nil(t, err)

	err = discovery.SyncAliasFile(&conf)
	assert.Nil(t, err)
}
