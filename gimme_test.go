package main

import (
	"fmt"
	"github.com/faceplate-kleo/gimme/config"
	"github.com/faceplate-kleo/gimme/discovery"
	"github.com/stretchr/testify/assert"
	"testing"
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
	paths, err := discovery.DiscoverPaths()
	assert.Nil(t, err)

	fmt.Println(paths)
}
