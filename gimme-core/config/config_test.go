package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestProcess(t *testing.T) {
	conf := Config{
		Dryrun:   false,
		Manifest: false,
	}
	data, err := ReadGimmeFileV1("../test/example-project/.gimme.yaml", &conf)
	if err != nil {
		t.Fatalf("Fatal: %s\n", err)
	}

	err = data.Process()
	if err != nil {
		t.Fatalf("Fatal: %s\n", err)
	}
}

func TestLoadConfig(t *testing.T) {
	fmt.Println(os.Getenv("PWD"))
	err := os.Setenv("GIMME_ROOT_PATH", "../test/.gimme/test-config.yaml")
	assert.Nil(t, err)
	conf, err := LoadRootConfig(false, false)
	assert.Nil(t, err)

	assert.Equal(t, "../test", conf.HomeOverride)
	assert.Equal(t, true, conf.KubeconfigEdit)
}
