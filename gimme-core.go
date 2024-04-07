package main

import (
	"flag"
	"fmt"
	"github.com/faceplate-kleo/gimme-core/config"
	"github.com/faceplate-kleo/gimme-core/discovery"
	"github.com/faceplate-kleo/gimme-core/warp"
	"log"
	"os"
)

func main() {
	var verboseFlag bool
	var dryrunFlag bool

	flag.BoolVar(&verboseFlag, "v", false, "toggle verbose output")
	flag.BoolVar(&dryrunFlag, "--dryrun", false, "toggle dry run (no file operations will be executed)")

	flag.Parse()

	conf, err := config.LoadRootConfig()
	if err != nil {
		log.Fatalf("fatal: %s", err)
	}

	for _, arg := range flag.Args() {
		if arg == "discover" {
			fmt.Printf("Discovering .gimme.yaml files starting from root %s...\n", conf.DetermineHome())
			paths, err := discovery.DiscoverPaths(conf.DetermineHome())
			if err != nil {
				log.Fatalf("fatal: %s", err)
			}

			if len(paths) == 0 {
				fmt.Println("No .gimme.yaml files found")
				os.Exit(1)
			}

			for _, path := range paths {
				fmt.Println(path)
			}
		} else if arg == "init" {
			err := conf.EnsureSetup()
			if err != nil {
				log.Fatalf("fatal during init: %s", err)
			}
		} else if arg == "sync" {
			err := discovery.SyncAliasFile(&conf)
			if err != nil {
				log.Fatalf("fatal during sync: %s", err)
			}
		} else {
			err := warp.Warp(arg, &conf)
			if err != nil {
				log.Fatalf("fatal during warp: %s", err)
			}
		}
	}
}
