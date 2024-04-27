package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type SpecV1 struct {
	Conf         *Config      `yaml:",omitempty"`
	GimmeVersion string       `yaml:"gimmeVersion"`
	Gimme        gimmeBlockV1 `yaml:"gimme"`
}

type gimmeBlockV1 struct {
	Alias string            `yaml:"alias"`
	Init  []string          `yaml:"init"`
	Env   map[string]string `yaml:"env"`
}

func (s *SpecV1) Process() error {
	err := s.setEnvironment()
	if err != nil {
		return err
	}
	err = s.executeInit()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpecV1) executeInit() error {
	if s.Conf.AutoTrust && s.Conf.Manifest {
		fmt.Println("[AUTOTRUST]")
	}
	for _, initLine := range s.Gimme.Init {
		if s.Conf.Manifest {
			fmt.Printf("[CMD] %s\n", initLine)
			continue
		}
		tokens := strings.Split(initLine, " ")
		args := make([]string, 0)
		if len(tokens) > 1 {
			args = tokens[1:]
		}
		cmd := exec.Command(tokens[0], args...)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("error executing init shell command \"%s\": %s", initLine, err)
		}
	}
	return nil
}

func (s *SpecV1) setEnvironment() error {
	for envKey, envVal := range s.Gimme.Env {
		if s.Conf.Manifest {
			fmt.Printf("[ENV] %s %s\n", envKey, envVal)
			continue
		}
		err := os.Setenv(envKey, envVal)
		if err != nil {
			return fmt.Errorf("error setting environment variable %s to %s : %s", envKey, envVal, err)
		}
	}
	return nil
}
