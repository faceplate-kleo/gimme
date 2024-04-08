package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestProcess(t *testing.T) {
	data, err := ReadGimmeFileV1("../test/example-project/.gimme.yaml")
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
	err := os.Setenv("GIMME_CONFIG_PATH", "../test/.gimme/test-config.yaml")
	assert.Nil(t, err)
	conf, err := LoadRootConfig()
	assert.Nil(t, err)

	assert.Equal(t, "../test", conf.HomeOverride)
	assert.Equal(t, true, conf.KubeconfigEdit)
}
