package main

import (
	"flag"
	"fmt"
	"github.com/faceplate-kleo/gimme-core/config"
	"github.com/faceplate-kleo/gimme-core/discovery"
	"github.com/faceplate-kleo/gimme-core/warp"
	"log"
	"os"
	"path"
)

func main() {
	var verboseFlag bool
	var dryrunFlag bool
	var manifestFlag bool

	flag.BoolVar(&verboseFlag, "v", false, "toggle verbose output")
	flag.BoolVar(&dryrunFlag, "dryrun", false, "toggle dry run (no file operations will be executed)")
	flag.BoolVar(&manifestFlag, "manifest", false, "toggle manifest mode (for use in injection scripting)")

	flag.Parse()

	for i, arg := range flag.Args() {
		if arg == "discover" {
			conf, err := config.LoadRootConfig(dryrunFlag, manifestFlag)
			if err != nil {
				log.Fatalf("fatal: %s", err)
			}
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
			conf := config.Config{}
			err := conf.EnsureSetup()
			if err != nil {
				log.Fatalf("fatal during init: %s", err)
			}
		} else if arg == "sync" {
			conf, err := config.LoadRootConfig(dryrunFlag, manifestFlag)
			if err != nil {
				log.Fatalf("fatal: %s", err)
			}
			err = discovery.SyncAliasFile(&conf)
			if err != nil {
				log.Fatalf("fatal during sync: %s", err)
			}
		} else if arg == "list" {
			conf, err := config.LoadRootConfig(dryrunFlag, manifestFlag)
			if err != nil {
				log.Fatalf("fatal: %s", err)
			}
			err = discovery.ListAliases(&conf)
			if err != nil {
				log.Fatalf("fatal during list: %s", err)
			}
		} else if arg == "pin" {
			conf, err := config.LoadRootConfig(dryrunFlag, manifestFlag)
			if err != nil {
				log.Fatalf("fatal: %s", err)
			}
			here := os.Getenv("PWD")
			alias := path.Dir(here) // just default to the name of the current directory
			if i != len(flag.Args())-1 {
				alias = flag.Args()[i+1]
			}
			err = discovery.PinDirectory(&conf, path.Join(here, ".gimme.yaml"), alias)
			if err != nil {
				log.Fatalf("failed to pin: %s", err)
			}
			break
		} else {
			conf, err := config.LoadRootConfig(dryrunFlag, manifestFlag)
			if err != nil {
				log.Fatalf("fatal: %s", err)
			}
			gimmePath, err := warp.Warp(arg, &conf)
			if err != nil {
				log.Fatalf("fatal during warp: %s", err)
			}

			projectConf, err := config.ReadGimmeFileV1(gimmePath, &conf)
			if err != nil {
				log.Fatalf("fatal: failed to read config file %s: %s", gimmePath, err)
			}
			err = projectConf.Process()
			if err != nil {
				log.Fatalf("fatal: failed to process: %s", err)
			}
		}
	}
}
