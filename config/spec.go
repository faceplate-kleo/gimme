package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Spec interface {
	SpecGeneric | SpecV1
	Process() error
}

type SpecGeneric struct {
	GimmeVersion string `yaml:"gimmeVersion"`
}

func (sg *SpecGeneric) Process() error {
	// intentional no-op
	return nil
}

func unpackGimmeFile(gimmePath string) (version string, data []byte, err error) {
	var yamlData SpecGeneric
	data, err = os.ReadFile(gimmePath)
	if err != nil {
		return "unknown", data, err
	}

	err = yaml.Unmarshal(data, &yamlData)
	if err != nil {
		return "unknown", data, err
	}

	return yamlData.GimmeVersion, data, nil
}

func ReadGimmeFileV1(gimmePath string) (SpecV1, error) {
	spec := SpecV1{}

	version, data, err := unpackGimmeFile(gimmePath)
	if err != nil {
		return spec, fmt.Errorf("failed to unpack gimme file: %s", err)
	}

	if version != "v1" {
		return spec, fmt.Errorf("unrecognized gimmeVersion: %s", version)
	}

	spec = SpecV1{}
	err = yaml.Unmarshal(data, &spec)
	if err != nil {
		return spec, fmt.Errorf("failed to unmarshal gimme data: %s", err)
	}

	return spec, nil
}
